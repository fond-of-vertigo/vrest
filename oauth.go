package vrest

import (
	"context"
	"errors"
	"net/url"
	"time"
)

var ErrOAuthTokenRequestFailed = errors.New("failed to get new oauth token")

// OAuthConfig is the configuration for an OAuth token request.
type OAuthConfig struct {
	URL          string
	GrantType    string
	Scope        string
	ClientID     string
	ClientSecret string
}

// OAuthToken is an OAuth token.
type OAuthToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	// ValidUntil is the time when the token is no longer valid.
	// It is calculated when a new token is received.
	ValidUntil time.Time `json:"-"`
}

type oauthTokenGetter struct {
	config OAuthConfig
	client *Client
}

// GetToken is a TokenGetter implementation that requests a new OAuth token.
// No synchronization or locking is required in this function.
func (o *oauthTokenGetter) GetToken(ctx context.Context, oldToken Token) (Token, error) {
	body := url.Values{}
	body.Add("grant_type", o.config.GrantType)
	body.Add("scope", o.config.Scope)
	body.Add("client_id", o.config.ClientID)
	body.Add("client_secret", o.config.ClientSecret)

	var token OAuthToken
	err := o.client.NewRequestWithContext(ctx).
		SetBaseURL(o.config.URL).
		SetContentType("application/x-www-form-urlencoded").
		SetBody(body.Encode()).
		SetTokenRequest().
		SetResponseBody(&token).
		DoPost("")
	if err != nil {
		return nil, errors.Join(ErrOAuthTokenRequestFailed, err)
	}

	safetyMarginMinutes := 5
	token.ValidUntil = time.Now().
		Add(time.Duration(token.ExpiresIn) * time.Second).
		Add(time.Duration(safetyMarginMinutes) * time.Minute * -1)
	return &token, nil
}

// Token returns the actual token.
func (t *OAuthToken) Token() string {
	if t == nil {
		return ""
	}
	return t.AccessToken
}

// NeedsRefresh returns true if the token needs to be refreshed.
// There is a safety margin of 5 minutes.
func (t *OAuthToken) NeedsRefresh() bool {
	if t == nil || t.AccessToken == "" {
		return true
	}
	return time.Now().After(t.ValidUntil)
}

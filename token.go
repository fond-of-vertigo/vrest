package vrest

import (
	"context"
	"sync/atomic"
)

// Token is an interface that represents any kind of token.
// Tokens are used with Bearer authentication.
type Token interface {
	// Token returns the actual token
	Token() string

	// NeedsRefresh returns true if the token needs to be refreshed.
	// This should take some safety margin into account.
	NeedsRefresh() bool
}

// TokenGetter is a function that returns a new token.
// No locking is required in the implementation of this function.
// vrest will take care of locking.
type TokenGetter interface {
	// GetToken returns a new or refreshed token.
	// If a previous token exists, it is passed as oldToken.
	GetToken(ctx context.Context, oldToken Token) (Token, error)
}

// getValidToken returns a valid token. If the current token is invalid, it will be refreshed.
// This function is thread-safe.
func (c *Client) getValidToken(ctx context.Context) (Token, error) {
	token := c.token.Load()
	if token != nil && !token.NeedsRefresh() {
		return token, nil
	}

	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	// double check after lock if another thread
	// already refreshed the token
	token = c.token.Load()
	if token != nil && !token.NeedsRefresh() {
		return token, nil
	}

	newToken, err := c.TokenGetter.GetToken(ctx, token)
	if err != nil {
		return nil, err
	}
	c.token.Store(newToken)
	return newToken, nil
}

type atomicToken struct {
	value atomic.Value
}

func (t *atomicToken) Load() Token {
	if val := t.value.Load(); val != nil {
		return val.(Token)
	}
	return nil
}

func (t *atomicToken) Store(token Token) {
	t.value.Store(token)
}

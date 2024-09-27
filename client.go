package vrest

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"log/slog"
	"net/http"
	"reflect"
	"sync"
	"time"
)

type (
	Doer     func(req *Request) error
	HTTPDoer func(req *Request) (*http.Response, error)
)

type Client struct {
	BaseURL string

	ResponseBodyLimit int64
	TraceBodies       bool

	ContentType   string
	Authorization string
	TokenGetter   TokenGetter

	ErrorType reflect.Type

	Overridable Overridables

	httpClient *http.Client
	traceMaker TraceMaker
	logger     *slog.Logger
	token      atomicToken
	tokenMutex sync.Mutex
}

type Overridables struct {
	Do Doer

	DoHTTPRequest HTTPDoer

	IsSuccess func(req *Request) bool

	JSONMarshal func(req *Request, v interface{}) ([]byte, error)

	JSONUnmarshal func(req *Request, data []byte, v interface{}) error

	XMLMarshal func(req *Request, v interface{}) ([]byte, error)

	XMLUnmarshal func(req *Request, data []byte, v interface{}) error
}

var (
	ErrResponseNotUnmarshaled = errors.New("response was not unmarshaled")
	ErrInvalidRequest         = errors.New("invalid request")
)

// New creates a new client with a default timeout of 0.
// This means that the client will not have a timeout.
func New() *Client {
	return NewWithTimeout(0)
}

// NewWithTimeout creates a new client with the given timeout.
// If the timeout is 0, the client will not have a timeout.
func NewWithTimeout(timeout time.Duration) *Client {
	return NewWithClient(&http.Client{
		Timeout: timeout,
	})
}

// NewWithClient creates a new client with the given http.Client.
// This allows for customization of the client's behavior.
// With the Overridables field, the default behavior can be overridden.
func NewWithClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient:  httpClient,
		logger:      slog.Default(),
		TraceBodies: true,

		Overridable: Overridables{
			Do:            Do,
			DoHTTPRequest: DoHTTPRequest,
			IsSuccess:     IsSuccess,
			JSONMarshal:   JSONMarshal,
			JSONUnmarshal: JSONUnmarshal,
			XMLMarshal:    XMLMarshal,
			XMLUnmarshal:  XMLUnmarshal,
		},
	}
}

// SetLogger sets the logger for the client.
func (c *Client) SetLogger(logger *slog.Logger) *Client {
	c.logger = logger
	return c
}

// SetTraceMaker sets the trace maker for the client.
// The trace maker is used to create traces for requests.
func (c *Client) SetTraceMaker(traceMaker TraceMaker) *Client {
	c.traceMaker = traceMaker
	return c
}

// SetTraceBodies sets whether the client should trace request and response bodies.
func (c *Client) SetTraceBodies(value bool) *Client {
	c.TraceBodies = value
	return c
}

// SetErrorBodyType extracts the type info from the value
// parameter and sets it as the error body type for the client.
// This is used when a IsSuccess() returns false to unmarshal
// the body to the given type.
// The unmarhsaled error will be part of the returned error.
func (c *Client) SetErrorBodyType(value error) *Client {
	c.ErrorType = typeOf(value)
	return c
}

// SetBaseURL sets the base URL for the client.
func (c *Client) SetBaseURL(baseURL string) *Client {
	c.BaseURL = baseURL
	return c
}

// SetContentTypeJSON sets the default content type
// for request bodies to application/json.
func (c *Client) SetContentTypeJSON() *Client {
	return c.SetContentType("application/json")
}

// SetContentTypeXML sets the default content type
// for request bodies to text/xml.
func (c *Client) SetContentTypeXML() *Client {
	return c.SetContentType("text/xml")
}

// SetContentType sets the default content type
// for request bodies to the given type.
func (c *Client) SetContentType(contentType string) *Client {
	c.ContentType = contentType
	return c
}

// SetBasicAuth sets the basic auth header for the client.
func (c *Client) SetBasicAuth(username, password string) *Client {
	return c.SetAuthorization("Basic " + encodeBasicAuth(username, password))
}

// SetBearerAuth sets the bearer auth header for the client.
func (c *Client) SetBearerAuth(token string) *Client {
	return c.SetAuthorization("Bearer " + token)
}

// SetOAuth sets the OAuth configuration for the client.
// This automatically sets the oauth token getter for the client.
func (c *Client) SetOAuth(cfg OAuthConfig) *Client {
	return c.SetTokenGetter(&oauthTokenGetter{
		config: cfg,
		client: c,
	})
}

// SetTokenGetter sets a custom token getter for the client.
// See the readme and examples for how to implement a custom token getter.
func (c *Client) SetTokenGetter(tokenGetter TokenGetter) *Client {
	c.TokenGetter = tokenGetter
	return c
}

// SetAuthorization sets the authorization header for the client.
func (c *Client) SetAuthorization(auth string) *Client {
	c.Authorization = auth
	return c
}

// SetResponseBodyLimit sets the response body limit for the client.
// If the response body is larger than the limit, it will be truncated.
// If the limit is 0, the response body will not be limited which can be
// used by an attacker to perform a DoS attack.
func (c *Client) SetResponseBodyLimit(limit int64) *Client {
	c.ResponseBodyLimit = limit
	return c
}

// JSONMarshal marshals the given value into JSON
// without escaping HTML characters.
func JSONMarshal(req *Request, v interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	return buffer.Bytes(), err
}

// JSONUnmarshal unmarshals the given data into the given value.
func JSONUnmarshal(req *Request, data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// XMLMarshal marshals the given value into XML.
func XMLMarshal(req *Request, v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

// XMLUnmarshal unmarshals the given data into the given value.
func XMLUnmarshal(req *Request, data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

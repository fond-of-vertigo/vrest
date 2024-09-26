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

var ErrResponseNotUnmarshaled = errors.New("response was not unmarshaled")

func New() *Client {
	return NewWithTimeout(0)
}

func NewWithTimeout(timeout time.Duration) *Client {
	return NewWithClient(&http.Client{
		Timeout: timeout,
	})
}

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

func (c *Client) SetLogger(logger *slog.Logger) *Client {
	c.logger = logger
	return c
}

func (c *Client) SetTraceMaker(traceMaker TraceMaker) *Client {
	c.traceMaker = traceMaker
	return c
}

func (c *Client) SetTraceBodies(value bool) *Client {
	c.TraceBodies = value
	return c
}

func (c *Client) SetErrorBodyType(value error) *Client {
	c.ErrorType = typeOf(value)
	return c
}

func (c *Client) SetBaseURL(baseURL string) *Client {
	c.BaseURL = baseURL
	return c
}

func (c *Client) SetContentTypeJSON() *Client {
	return c.SetContentType("application/json")
}

func (c *Client) SetContentTypeXML() *Client {
	return c.SetContentType("text/xml")
}

func (c *Client) SetContentType(contentType string) *Client {
	c.ContentType = contentType
	return c
}

func (c *Client) SetBasicAuth(username, password string) *Client {
	return c.SetAuthorization("Basic " + encodeBasicAuth(username, password))
}

func (c *Client) SetBearerAuth(token string) *Client {
	return c.SetAuthorization("Bearer " + token)
}

func (c *Client) SetOAuth(cfg OAuthConfig) *Client {
	return c.SetTokenGetter(&oauthTokenGetter{
		config: cfg,
		client: c,
	})
}

func (c *Client) SetTokenGetter(tokenGetter TokenGetter) *Client {
	c.TokenGetter = tokenGetter
	return c
}

func (c *Client) SetAuthorization(auth string) *Client {
	c.Authorization = auth
	return c
}

func (c *Client) SetResponseBodyLimit(limit int64) *Client {
	c.ResponseBodyLimit = limit
	return c
}

func JSONMarshal(req *Request, v interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	return buffer.Bytes(), err
}

func JSONUnmarshal(req *Request, data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func XMLMarshal(req *Request, v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func XMLUnmarshal(req *Request, data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

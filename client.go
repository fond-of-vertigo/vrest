package vrest

import (
	"encoding/json"
	"encoding/xml"
	"log/slog"
	"net/http"
	"reflect"
	"time"
)

type Doer func(req *Request) error
type HTTPDoer func(req *Request) (*http.Response, error)

type Client struct {
	BaseURL string

	ResponseBodyLimit int64

	ContentType   string
	Authorization string

	ErrorType reflect.Type

	Overridable Overridables

	httpClient *http.Client
	traceMaker TraceMaker
	logger     *slog.Logger
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
		httpClient: httpClient,
		logger:     slog.Default(),

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

func (c *Client) SetAuthorization(auth string) *Client {
	c.Authorization = auth
	return c
}

func (c *Client) SetResponseBodyLimit(limit int64) *Client {
	c.ResponseBodyLimit = limit
	return c
}

func JSONMarshal(req *Request, v interface{}) ([]byte, error) {
	return json.Marshal(v)
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

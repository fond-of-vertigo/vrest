package vrest

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"reflect"
	"time"
)

type Doer func(req *Request) error
type HTTPDoer func(req *Request) (*http.Response, error)

type Client struct {
	BaseURL string

	ResponseBodyLimit int64

	ContentType string

	ErrorType reflect.Type

	Overridable Overridables

	httpClient *http.Client
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

// New creates a new client with http.DefaultClient
func New() *Client {
	return NewWithClient(http.DefaultClient)
}

func NewWithClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,

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

func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.httpClient.Timeout = timeout
	return c
}

func (c *Client) SetErrorBodyType(value interface{}) *Client {
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

func (c *Client) SetResponseBodyLimit(limit int64) *Client {
	c.ResponseBodyLimit = limit
	return c
}

func Do(req *Request) error {
	err := req.makeHTTPRequest()
	if err != nil {
		return err
	}

	rawResp, err := req.Overridable.DoHTTPRequest(req)
	defer req.Client.closeRawResponse(rawResp)

	err = req.processHTTPResponse(rawResp, err)
	return err
}

func DoHTTPRequest(req *Request) (*http.Response, error) {
	return req.Client.httpClient.Do(req.Raw)
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

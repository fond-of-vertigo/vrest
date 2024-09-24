package vrest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	Client        *Client
	Context       context.Context
	Raw           *http.Request
	Method        string
	Path          string
	Header        http.Header
	Query         url.Values
	Body          interface{}
	BodyBytes     []byte
	ContentLength int64
	Response      Response
	Overridable   Overridables
	TraceBody     bool
}

// NewRequest is a shortcut for NewRequestWithContext(context.Background()).
func (c *Client) NewRequest() *Request {
	return c.NewRequestWithContext(context.Background())
}

// NewRequestWithContext creates a new Request instance and initializes its fields with default values.
// It also sets the Content-Type and Authorization headers of the request if they are provided in the Client struct.
func (c *Client) NewRequestWithContext(ctx context.Context) *Request {
	req := &Request{
		Client:      c,
		Context:     ctx,
		Header:      make(http.Header),
		Query:       make(url.Values),
		Overridable: c.Overridable,
		TraceBody:   c.TraceBodies,
		Response: Response{
			BodyLimit:   c.ResponseBodyLimit,
			TraceBody:   c.TraceBodies,
			DoUnmarshal: true,
		},
	}

	if c.ContentType != "" {
		req.SetContentType(c.ContentType)
	}
	if c.Authorization != "" {
		req.SetAuthorization(c.Authorization)
	}

	return req
}

func (req *Request) makeHTTPRequest() error {
	if req.Raw != nil {
		return nil
	}

	reqURL := req.makeRequestURL(req.Client.BaseURL, req.Path)

	reqBodyReader, bodyBytes, err := req.makeRequestBody(req.Body, req.ContentType())
	if err != nil {
		return err
	}
	req.BodyBytes = bodyBytes

	req.Raw, err = http.NewRequestWithContext(req.Context, req.Method, reqURL, reqBodyReader)
	if err != nil {
		return err
	}
	if req.ContentLength > 0 {
		req.Raw.ContentLength = req.ContentLength
	}

	if len(req.Header) > 0 {
		req.Raw.Header = req.Header
		if len(bodyBytes) == 0 && !req.bodyIsReader() {
			req.Raw.Header.Del("Content-Type")
		}
	}

	if len(req.Query) > 0 {
		req.Raw.URL.RawQuery = req.Query.Encode()
	}

	return nil
}

func (req *Request) makeRequestBody(body interface{}, contentType string) (io.Reader, []byte, error) {
	switch bodyValue := body.(type) {
	case io.Reader:
		return bodyValue, nil, nil
	case []byte:
		return bytes.NewReader(bodyValue), bodyValue, nil
	case string:
		bodyBytes := []byte(bodyValue)
		return bytes.NewReader(bodyBytes), bodyBytes, nil
	case nil:
		return nil, nil, nil
	default:
		bodyBytes, err := req.marshalRequestBody(body, contentType)
		if err != nil {
			return nil, nil, err
		}
		return bytes.NewReader(bodyBytes), bodyBytes, nil
	}
}

func (req *Request) bodyIsReader() bool {
	if r, ok := req.Body.(io.Reader); ok && r != nil {
		return true
	}
	return false
}

func (req *Request) marshalRequestBody(body interface{}, contentType string) ([]byte, error) {
	if contentType == "" {
		return nil, fmt.Errorf("content type must not be empty")
	}

	if IsJSONContentType(contentType) {
		return req.Overridable.JSONMarshal(req, body)
	}

	if IsXMLContentType(contentType) {
		return req.Overridable.XMLMarshal(req, body)
	}

	return nil, fmt.Errorf("don't know how to marshal request body with Content-Type \"%s\"", contentType)
}

func (req *Request) makeRequestURL(baseURL, requestURL string) string {
	if baseURL != "" && !strings.Contains(requestURL, baseURL) {
		return baseURL + requestURL
	}
	return requestURL
}

// SetContext sets the context of the request.
func (req *Request) SetContext(ctx context.Context) *Request {
	req.Context = ctx
	return req
}

// SetBody sets the body of the request.
func (req *Request) SetBody(body interface{}) *Request {
	req.Body = body
	return req
}

// SetTraceRequestBody sets the TraceBody field of the request.
func (req *Request) SetTraceRequestBody(value bool) *Request {
	req.TraceBody = value
	return req
}

// SetTraceResponseBody sets the TraceBody field of the response.
func (req *Request) SetTraceResponseBody(value bool) *Request {
	req.Response.TraceBody = value
	return req
}

// SetQueryParamIf sets the query parameter of the request when the condition matches.
func (req *Request) SetQueryParamIf(condition bool, key string, values ...string) *Request {
	if condition {
		return req.SetQueryParam(key, values...)
	}
	return req
}

// SetQueryParam sets the query parameter of the request.
func (req *Request) SetQueryParam(key string, values ...string) *Request {
	req.Query[key] = values
	return req
}

// ContentType returns the Content-Type header of the request.
func (req *Request) ContentType() string {
	return req.Header.Get("Content-Type")
}

// SetContentLength sets the Content-Length header of the request.
func (req *Request) SetContentLength(contentLength int64) *Request {
	req.ContentLength = contentLength
	return req
}

// SetResponseBody sets the body of the response.
func (req *Request) SetResponseBody(value interface{}) *Request {
	req.Response.Body = value
	return req
}

// SetResponseBodyLimit sets the limit of the response body.
func (req *Request) SetResponseBodyLimit(limit int64) *Request {
	req.Response.BodyLimit = limit
	return req
}

// SetResponseContentLengthPtr sets the Content-Length based on a pointer.
func (req *Request) SetResponseContentLengthPtr(contentLengthPtr *int64) *Request {
	req.Response.ContentLengthPtr = contentLengthPtr
	return req
}

// SetSuccessStatusCode sets the success status code of the response.
func (req *Request) SetSuccessStatusCode(statusCodes ...int) *Request {
	req.Response.SuccessStatusCodes = statusCodes
	return req
}

// ForceResponseJSON forces the response to be JSON.
// This is useful when the server does not return the correct Content-Type header.
func (req *Request) ForceResponseJSON() *Request {
	req.Response.ForceJSON = true
	return req
}

// ForceResponseXML forces the response to be XML.
// This is useful when the server does not return the correct Content-Type header.
func (req *Request) ForceResponseXML() *Request {
	req.Response.ForceXML = true
	return req
}

// SetResponseErrorBody sets the error body of the response.
func (req *Request) SetResponseErrorBody(value interface{}) *Request {
	req.Response.ErrorBody = value
	return req
}

// SetBasicAuth sets the basic authentication of the request.
// The username and password are separated by a colon and then base64 encoded.
func (req *Request) SetBasicAuth(username, password string) *Request {
	return req.SetAuthorization("Basic " + encodeBasicAuth(username, password))
}

// SetBearerAuth sets the bearer authentication of the request.
// The token should not contain the "Bearer " prefix.
func (req *Request) SetBearerAuth(token string) *Request {
	return req.SetAuthorization("Bearer " + token)
}

// SetContentTypeJSON sets the Content-Type header of the request to "application/json".
func (req *Request) SetContentTypeJSON() *Request {
	return req.SetContentType("application/json")
}

// SetContentTypeXML sets the Content-Type header of the request to "text/xml".
func (req *Request) SetContentTypeXML() *Request {
	return req.SetContentType("text/xml")
}

// SetContentType sets a custom Content-Type header of the request.
func (req *Request) SetContentType(contentType string) *Request {
	return req.SetHeader("Content-Type", contentType)
}

// SetAuthorization sets the Authorization header of the request to a custom authentication string.
func (req *Request) SetAuthorization(authValue string) *Request {
	return req.SetHeader("Authorization", authValue)
}

// SetHeaderIf sets the header of the request when the condition matches.
func (req *Request) SetHeaderIf(condition bool, key string, value string) *Request {
	if condition {
		return req.SetHeader(key, value)
	}
	return req
}

// SetHeader sets the header of the request.
func (req *Request) SetHeader(key string, value string) *Request {
	req.Header.Set(key, value)
	return req
}

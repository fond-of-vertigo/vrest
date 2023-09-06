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
	Client      *Client
	Context     context.Context
	Raw         *http.Request
	Method      string
	Path        string
	Header      http.Header
	Query       url.Values
	Body        interface{}
	BodyBytes   []byte
	Response    Response
	Overridable Overridables
}

func (c *Client) NewRequest() *Request {
	req := &Request{
		Client:      c,
		Context:     context.Background(),
		Header:      make(http.Header),
		Query:       make(url.Values),
		Overridable: c.Overridable,
		Response: Response{
			BodyLimit: c.ResponseBodyLimit,
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

	if len(req.Header) > 0 {
		req.Raw.Header = req.Header
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

func (req *Request) SetContext(ctx context.Context) *Request {
	req.Context = ctx
	return req
}

func (req *Request) SetBody(body interface{}) *Request {
	req.Body = body
	return req
}

func (req *Request) SetQueryParamIf(condition bool, key string, values ...string) *Request {
	if condition {
		return req.SetQueryParam(key, values...)
	}
	return req
}

func (req *Request) SetQueryParam(key string, values ...string) *Request {
	req.Query[key] = values
	return req
}

func (req *Request) ContentType() string {
	return req.Header.Get("Content-Type")
}

func (req *Request) SetResponseBody(value interface{}) *Request {
	req.Response.Body = value
	return req
}

func (req *Request) SetResponseBodyLimit(limit int64) *Request {
	req.Response.BodyLimit = limit
	return req
}

func (req *Request) ForceResponseJSON() *Request {
	req.Response.ForceJSON = true
	return req
}

func (req *Request) ForceResponseXML() *Request {
	req.Response.ForceXML = true
	return req
}

func (req *Request) SetResponseErrorBody(value interface{}) *Request {
	req.Response.ErrorBody = value
	return req
}

func (req *Request) SetBasicAuth(username, password string) *Request {
	return req.SetAuthorization("Basic " + encodeBasicAuth(username, password))
}

func (req *Request) SetBearerAuth(token string) *Request {
	return req.SetAuthorization("Bearer " + token)
}

func (req *Request) SetContentTypeJSON() *Request {
	return req.SetContentType("application/json")
}

func (req *Request) SetContentTypeXML() *Request {
	return req.SetContentType("text/xml")
}

func (req *Request) SetContentType(contentType string) *Request {
	return req.SetHeader("Content-Type", contentType)
}

func (req *Request) SetAuthorization(authValue string) *Request {
	return req.SetHeader("Authorization", authValue)
}

func (req *Request) SetHeaderIf(condition bool, key string, value string) *Request {
	if condition {
		return req.SetHeader(key, value)
	}
	return req
}

func (req *Request) SetHeader(key string, value string) *Request {
	req.Header.Set(key, value)
	return req
}

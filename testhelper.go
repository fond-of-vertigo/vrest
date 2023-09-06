package vrest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

func MockDoer(responseValue interface{}, err error) Doer {
	return func(req *Request) error {
		return mockSetResponseValue(req, responseValue)
	}
}

type MockHTTPParams struct {
	StatusCode  int
	Body        []byte
	BodyString  string
	BodyReader  io.Reader
	ContentType string
	Error       error
}

func MockHTTPDoer(p MockHTTPParams, additionalHeaders ...string) HTTPDoer {
	if len(additionalHeaders)%2 != 0 {
		panic("len(additionalHeaders) is not even!")
	}

	resp := http.Response{
		StatusCode: p.StatusCode,
		Header:     make(http.Header),
	}

	if resp.StatusCode == 0 {
		resp.StatusCode = 200
	}

	if p.BodyReader != nil {
		resp.Body = io.NopCloser(p.BodyReader)
	} else if p.BodyString != "" {
		resp.Body = io.NopCloser(bytes.NewReader([]byte(p.BodyString)))
	} else {
		resp.Body = io.NopCloser(bytes.NewReader(p.Body))
	}

	if p.ContentType != "" {
		resp.Header.Set("Content-Type", p.ContentType)
	}

	for i := 0; i < len(additionalHeaders); i += 2 {
		resp.Header.Set(additionalHeaders[i], additionalHeaders[i+1])
	}

	return func(req *Request) (*http.Response, error) {
		return &resp, p.Error
	}
}

func MockJSONResponse(statusCode int, body string) MockHTTPParams {
	return MockHTTPParams{
		StatusCode:  statusCode,
		BodyString:  body,
		ContentType: "application/json",
	}
}

func MockXMLResponse(statusCode int, body string) MockHTTPParams {
	return MockHTTPParams{
		StatusCode:  statusCode,
		BodyString:  body,
		ContentType: "text/xml",
	}
}

func mockSetResponseValue(req *Request, value interface{}) error {
	if req.Response.Body == nil {
		return fmt.Errorf("req.Response.Body is nil")
	}
	if value == nil {
		return nil
	}

	target := reflect.ValueOf(req.Response.Body)
	if target.Kind() != reflect.Ptr {
		return fmt.Errorf("req.Response.Body must be a pointer")
	}

	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	target.Elem().Set(val)
	return nil
}

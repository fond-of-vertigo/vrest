package vrest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func MockDoer(responseValue interface{}, err error) Doer {
	return func(req *Request) error {
		return mockSetResponseValue(req, responseValue)
	}
}

type MockHTTPResponse struct {
	StatusCode  int
	Body        []byte
	BodyString  string
	BodyReader  io.Reader
	ContentType string
	Error       error

	CapturedRequest *Request
}

func MockHTTPDoer(p *MockHTTPResponse, additionalHeaders ...string) HTTPDoer {
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
		p.CapturedRequest = req
		return &resp, p.Error
	}
}

func MockJSONResponseFromFile(t *testing.T, statusCode int, filePath string) *MockHTTPResponse {
	return &MockHTTPResponse{
		StatusCode:  statusCode,
		Body:        MustReadFile(t, filePath),
		ContentType: "application/json",
	}
}

func MockJSONResponse(statusCode int, body string) *MockHTTPResponse {
	return &MockHTTPResponse{
		StatusCode:  statusCode,
		BodyString:  body,
		ContentType: "application/json",
	}
}

func MockXMLResponseFromFile(t *testing.T, statusCode int, filePath string) *MockHTTPResponse {
	return &MockHTTPResponse{
		StatusCode:  statusCode,
		Body:        MustReadFile(t, filePath),
		ContentType: "text/xml",
	}
}

func MockXMLResponse(statusCode int, body string) *MockHTTPResponse {
	return &MockHTTPResponse{
		StatusCode:  statusCode,
		BodyString:  body,
		ContentType: "text/xml",
	}
}

func MustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	bytes, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatal(err)
	}
	return bytes
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

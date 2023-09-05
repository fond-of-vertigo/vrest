package vrest

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
)

type Response struct {
	Raw       *http.Response
	Body      interface{}
	ErrorBody interface{}
	ForceJSON bool
	ForceXML  bool
	BodyBytes []byte
	BodyLimit int64
}

func (req *Request) processHTTPResponse(rawResp *http.Response, err error) error {
	req.Response.Raw = rawResp
	if err != nil {
		return fmt.Errorf("http request %s %s failed: %w", req.Raw.Method, req.Raw.URL, err)
	}
	if rawResp == nil {
		return fmt.Errorf("http request %s %s returned no response and no error", req.Raw.Method, req.Raw.URL)
	}

	err = req.readResponseBody()
	if err != nil {
		return fmt.Errorf("http request %s %s failed to read response body: %w", req.Raw.Method, req.Raw.URL, err)
	}

	success := req.Overridable.IsSuccess(req)
	if req.Response.HasEmptyBody() {
		if !success {
			return fmt.Errorf("http request %s %s failed with status code %d", req.Raw.Method, req.Raw.URL, req.Response.StatusCode())
		}
		return nil
	}

	responseValue := req.Response.Body
	if !success {
		if req.Response.ErrorBody == nil && req.Client.ErrorType != nil {
			req.Response.ErrorBody = reflect.New(req.Client.ErrorType).Interface()
		}
		responseValue = req.Response.ErrorBody
	}

	didUnmarshal, err := req.unmarshalResponseBody(responseValue)
	if err != nil {
		return fmt.Errorf("http request %s %s failed to unmarshal response body: %w", req.Raw.Method, req.Raw.URL, err)
	}

	if !success {
		if didUnmarshal {
			switch e := responseValue.(type) {
			case error:
				return fmt.Errorf("http request %s %s failed: %w", req.Raw.Method, req.Raw.URL, e)
			default:
				return fmt.Errorf("http request %s %s failed: %s", req.Raw.Method, req.Raw.URL, responseValue)
			}
		} else {
			msg := string(req.Response.BodyBytes)
			return fmt.Errorf("http request %s %s failed: %s", req.Raw.Method, req.Raw.URL, msg)
		}
	}

	return nil
}

func (req *Request) readResponseBody() error {
	if req.Response.Raw.Body == nil {
		return nil
	}

	var r io.Reader = req.Response.Raw.Body
	if req.Response.BodyLimit > 0 {
		r = io.LimitReader(r, req.Response.BodyLimit)
	}

	var err error
	req.Response.BodyBytes, err = io.ReadAll(r)
	return err
}

func (req *Request) unmarshalResponseBody(value interface{}) (bool, error) {
	var err error
	if req.Response.ForceJSON || IsJSONContentType(req.Response.ContentType()) {
		err = req.Overridable.JSONUnmarshal(req, req.Response.BodyBytes, value)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	if req.Response.ForceXML || IsXMLContentType(req.Response.ContentType()) {
		err = req.Overridable.XMLUnmarshal(req, req.Response.BodyBytes, value)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func (resp *Response) HasEmptyBody() bool {
	return resp.Raw.Body == nil
}

func (resp *Response) StatusCode() int {
	return resp.Raw.StatusCode
}

func (resp *Response) ContentType() string {
	return resp.Header().Get("Content-Type")
}

func (resp *Response) Header() http.Header {
	return resp.Raw.Header
}

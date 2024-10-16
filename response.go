package vrest

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
)

type Response struct {
	Raw         *http.Response
	Error       error
	Body        any
	ErrorBody   any
	ForceJSON   bool
	ForceXML    bool
	BodyBytes   []byte
	BodyLimit   int64
	TraceBody   bool
	DoUnmarshal bool

	ContentLengthPtr *int64

	SuccessStatusCodes []int
}

// processHTTPResponse processes the raw http response and sets the response fields.
// It returns an error if the response is not successful or if the response body could not be read.
func (req *Request) processHTTPResponse(rawResp *http.Response, err error) error {
	req.Response.Raw = rawResp
	if err != nil {
		return fmt.Errorf("http request %s %s failed: %w", req.Raw.Method, req.Raw.URL, err)
	}
	if rawResp == nil {
		return fmt.Errorf("http request %s %s returned no response and no error", req.Raw.Method, req.Raw.URL)
	}

	if req.Response.ContentLengthPtr != nil {
		*req.Response.ContentLengthPtr = req.Response.Raw.ContentLength
	}

	err = req.readResponseBody()
	if err != nil {
		return fmt.Errorf("http request %s %s failed to read response body: %w", req.Raw.Method, req.Raw.URL, err)
	}

	success := req.Overridable.IsSuccess(req)
	if req.Response.HasEmptyBody() {
		if !success {
			return fmt.Errorf("http request %s %s failed with status code %d",
				req.Raw.Method, req.Raw.URL, req.Response.StatusCode())
		}
		return nil
	}

	responseValue := req.Response.Body

	// if the response is not successful, unmarshal the error body
	if !success && req.Response.DoUnmarshal {
		if req.Response.ErrorBody == nil && req.Client.ErrorType != nil {
			req.Response.ErrorBody = reflect.New(req.Client.ErrorType).Interface()
		}
		responseValue = req.Response.ErrorBody
	}

	didUnmarshal, err := req.unmarshalResponseBody(responseValue)
	if success && err != nil {
		// treat unmarshaling error as a failure only if the response is successful
		return fmt.Errorf("http request %s %s failed to unmarshal response body: %w", req.Raw.Method, req.Raw.URL, err)
	}

	if !success {
		if didUnmarshal {
			errMsg := fmt.Sprintf("http request %s %s failed: status %d", req.Raw.Method, req.Raw.URL, req.Response.StatusCode())
			switch e := responseValue.(type) {
			case error:
				return fmt.Errorf("%s: %w", errMsg, e)
			default:
				return fmt.Errorf("%s: %s", errMsg, responseValue)
			}
		}

		msg := string(req.Response.BodyBytes)
		return fmt.Errorf("http request %s %s failed: status %d: %s", req.Raw.Method, req.Raw.URL,
			req.Response.StatusCode(), msg)
	}

	return nil
}

// readResponseBody reads the response body and sets the response body bytes.
// It returns an error if the response body could not be read.
// It does not read the body if the response body is of type io.ReadCloser.
// It checks for a response body limit and reads only up to that limit, if
// request.SetResponseBodyLimit was called.
func (req *Request) readResponseBody() error {
	if req.Response.Raw.Body == nil {
		return nil
	}

	if req.Overridable.IsSuccess(req) {
		// check if we should just return the response ReadCloser
		if responseReadCloser, ok := req.Response.Body.(*io.ReadCloser); ok && responseReadCloser != nil {
			*responseReadCloser = req.Response.Raw.Body
			req.Response.DoUnmarshal = false
			return nil
		}
	}

	var r io.Reader = req.Response.Raw.Body
	if req.Response.BodyLimit > 0 {
		r = io.LimitReader(r, req.Response.BodyLimit)
	}

	var err error
	req.Response.BodyBytes, err = io.ReadAll(r)
	if len(req.Response.BodyBytes) > 0 && req.Response.WantsRawByteArray() {
		if responseBytesPointer, ok := req.Response.Body.(*[]byte); ok && responseBytesPointer != nil {
			*responseBytesPointer = req.Response.BodyBytes
			req.Response.DoUnmarshal = false
		}
	}

	return err
}

// unmarshalResponseBody unmarshals the response body into the given value.
// It returns true if the response body was unmarshaled, false otherwise.
// If there is a body and it was not unmarshaled, an error is returned.
func (req *Request) unmarshalResponseBody(value interface{}) (bool, error) {
	if !req.Response.DoUnmarshal {
		return false, nil
	}

	if value == nil {
		// seems like the caller is not interested in the result value
		// just exit without error
		return false, nil
	}

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

	if len(req.Response.BodyBytes) > 0 {
		return false, ErrResponseNotUnmarshaled
	}

	return false, nil
}

// WantsRawByteArray returns whether the response wants a raw byte array.
func (resp *Response) WantsRawByteArray() bool {
	if resp == nil || resp.Body == nil {
		return false
	}
	if responseBytesPointer, ok := resp.Body.(*[]byte); ok && responseBytesPointer != nil {
		return true
	}
	return false
}

// WantsReadCloser returns whether the response wants a ReadCloser
// which is useful for streaming the response body.
func (resp *Response) WantsReadCloser() bool {
	if resp == nil || resp.Body == nil {
		return false
	}
	if responseReadCloserPointer, ok := resp.Body.(*io.ReadCloser); ok && responseReadCloserPointer != nil {
		return true
	}
	return false
}

// HasEmptyBody returns whether the response has an empty body.
func (resp *Response) HasEmptyBody() bool {
	if resp.Raw == nil {
		return false
	}
	return resp.Raw.Body == nil
}

// StatusCode returns the status code of the response.
func (resp *Response) StatusCode() int {
	if resp.Raw == nil {
		return 0
	}
	return resp.Raw.StatusCode
}

// ContentType returns the Content-Type header of the response.
func (resp *Response) ContentType() string {
	return resp.Header().Get("Content-Type")
}

// Header returns the header map of the response.
func (resp *Response) Header() http.Header {
	if resp.Raw == nil {
		return http.Header{}
	}
	return resp.Raw.Header
}

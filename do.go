package vrest

import (
	"fmt"
	"net/http"
	"reflect"
)

// Do sends the request.
// If the client has a trace maker set, it will create a trace.
func Do(req *Request) error {
	err := req.makeHTTPRequest()
	if err != nil {
		return err
	}

	var trace Trace
	if req.Client.traceMaker != nil {
		trace = req.Client.traceMaker.NewTrace(req)
		defer trace.End()
	}

	req.Response.Raw, err = req.Overridable.DoHTTPRequest(req)
	if req.shouldCloseResponseBody() {
		defer req.Client.closeRawResponse(req)
	}

	// Processing the HTTP response.
	// This will read the response body and unmarshal it if needed.
	// It will also check if the response is successful.
	// If the response is not successful, it will return an error.
	req.Response.Error = req.processHTTPResponse(req.Response.Raw, err)

	if trace != nil {
		trace.OnAfterRequest(req)
	}

	return req.Response.Error
}

// DoGet sends the request with the GET method.
// It only works with static paths.
// If you need to use dynamic paths, use DoGetf.
func (req *Request) DoGet(path string) error {
	return req.Do(http.MethodGet, path)
}

// DoGetf sends the request with the GET method.
// It works with dynamic paths.
// The pathFormat is a format string that will be used with fmt.Sprintf.
func (req *Request) DoGetf(pathFormat string, values ...any) error {
	return req.Dof(http.MethodGet, pathFormat, values...)
}

// DoHead sends the request with the HEAD method.
// It only works with static paths.
// If you need to use dynamic paths, use DoHeadf.
func (req *Request) DoHead(path string) error {
	return req.Do(http.MethodHead, path)
}

// DoHeadf sends the request with the HEAD method.
// It works with dynamic paths.
// The pathFormat is a format string that will be used with fmt.Sprintf.
func (req *Request) DoHeadf(pathFormat string, values ...any) error {
	return req.Dof(http.MethodHead, pathFormat, values...)
}

// DoPost sends the request with the POST method.
// It only works with static paths.
// If you need to use dynamic paths, use DoPostf.
func (req *Request) DoPost(path string) error {
	return req.Do(http.MethodPost, path)
}

// DoPostf sends the request with the POST method.
// It works with dynamic paths.
// The pathFormat is a format string that will be used with fmt.Sprintf.
func (req *Request) DoPostf(pathFormat string, values ...any) error {
	return req.Dof(http.MethodPost, pathFormat, values...)
}

// DoPut sends the request with the PUT method.
// It only works with static paths.
// If you need to use dynamic paths, use DoPutf.
func (req *Request) DoPut(path string) error {
	return req.Do(http.MethodPut, path)
}

// DoPutf sends the request with the PUT method.
// It works with dynamic paths.
// The pathFormat is a format string that will be used with fmt.Sprintf.
func (req *Request) DoPutf(pathFormat string, values ...any) error {
	return req.Dof(http.MethodPut, pathFormat, values...)
}

// DoPatch sends the request with the PATCH method.
// It only works with static paths.
// If you need to use dynamic paths, use DoPatchf.
func (req *Request) DoPatch(path string) error {
	return req.Do(http.MethodPatch, path)
}

// DoPatchf sends the request with the PATCH method.
// It works with dynamic paths.
// The pathFormat is a format string that will be used with fmt.Sprintf.
func (req *Request) DoPatchf(pathFormat string, values ...any) error {
	return req.Dof(http.MethodPatch, pathFormat, values...)
}

// DoDelete sends the request with the DELETE method.
// It only works with static paths.
// If you need to use dynamic paths, use DoDeletef.
func (req *Request) DoDelete(path string) error {
	return req.Do(http.MethodDelete, path)
}

// DoDeletef sends the request with the DELETE method.
// It works with dynamic paths.
// The pathFormat is a format string that will be used with fmt.Sprintf.
func (req *Request) DoDeletef(pathFormat string, values ...any) error {
	return req.Dof(http.MethodDelete, pathFormat, values...)
}

// DoConnect sends the request with the CONNECT method.
// It only works with static paths.
// If you need to use dynamic paths, use DoConnectf.
func (req *Request) DoConnect(path string) error {
	return req.Do(http.MethodConnect, path)
}

// DoConnectf sends the request with the CONNECT method.
// It works with dynamic paths.
// The pathFormat is a format string that will be used with fmt.Sprintf.
func (req *Request) DoConnectf(pathFormat string, values ...any) error {
	return req.Dof(http.MethodConnect, pathFormat, values...)
}

// DoOptions sends the request with the OPTIONS method.
// It only works with static paths.
// If you need to use dynamic paths, use DoOptionsf.
func (req *Request) DoOptions(path string) error {
	return req.Do(http.MethodOptions, path)
}

// DoOptionsf sends the request with the OPTIONS method.
// It works with dynamic paths.
// The pathFormat is a format string that will be used with fmt.Sprintf.
func (req *Request) DoOptionsf(pathFormat string, values ...any) error {
	return req.Dof(http.MethodOptions, pathFormat, values...)
}

// DoTrace sends the request with the TRACE method.
// It only works with static paths.
// If you need to use dynamic paths, use DoTracef.
func (req *Request) DoTrace(path string) error {
	return req.Do(http.MethodTrace, path)
}

// DoTracef sends the request with the TRACE method.
// It works with dynamic paths.
// The pathFormat is a format string that will be used with fmt.Sprintf.
func (req *Request) DoTracef(pathFormat string, values ...any) error {
	return req.Dof(http.MethodTrace, pathFormat, values...)
}

// Dof sends the request with the given method. It works with dynamic paths.
func (req *Request) Dof(method, pathFormat string, values ...any) error {
	path := fmt.Sprintf(pathFormat, values...)
	return req.Do(method, path)
}

// Do sends the request with the given method.
func (req *Request) Do(method, path string) error {
	req.Method = method
	req.Path = path

	if err := req.validateBeforeDo(); err != nil {
		return err
	}

	return req.Client.Overridable.Do(req)
}

// DoHTTPRequest sends the request using the http.Client.
func DoHTTPRequest(req *Request) (*http.Response, error) {
	return req.Client.httpClient.Do(req.Raw)
}

func (req *Request) shouldCloseResponseBody() bool {
	if req.Overridable.IsSuccess(req) && req.Response.WantsReadCloser() {
		// When the caller wants a ReadCloser as result, we don't close
		// the response body for the caller.
		// But only for requests that were sucessful, so the caller still
		// gets the error unmarshalling for free.
		return false
	}
	return true
}

func (req *Request) validateBeforeDo() error {
	if req.Response.Body != nil && reflect.ValueOf(req.Response.Body).Kind() != reflect.Ptr {
		return fmt.Errorf("%w: the value you passed to request.SetResponseBody() must be a pointer", ErrInvalidRequest)
	}
	return nil
}

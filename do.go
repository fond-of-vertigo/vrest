package vrest

import (
	"fmt"
	"net/http"
)

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

	req.Response.Error = req.processHTTPResponse(req.Response.Raw, err)

	if trace != nil {
		trace.OnAfterRequest(req)
	}

	return req.Response.Error
}

func (req *Request) DoGet(path string) error {
	return req.Do(http.MethodGet, path)
}

func (req *Request) DoGetf(pathFormat string, values ...any) error {
	return req.Dof(http.MethodGet, pathFormat, values...)
}

func (req *Request) DoHead(path string) error {
	return req.Do(http.MethodHead, path)
}

func (req *Request) DoHeadf(pathFormat string, values ...any) error {
	return req.Dof(http.MethodHead, pathFormat, values...)
}

func (req *Request) DoPost(path string) error {
	return req.Do(http.MethodPost, path)
}

func (req *Request) DoPostf(pathFormat string, values ...any) error {
	return req.Dof(http.MethodPost, pathFormat, values...)
}

func (req *Request) DoPut(path string) error {
	return req.Do(http.MethodPut, path)
}

func (req *Request) DoPutf(pathFormat string, values ...any) error {
	return req.Dof(http.MethodPut, pathFormat, values...)
}

func (req *Request) DoPatch(path string) error {
	return req.Do(http.MethodPatch, path)
}

func (req *Request) DoPatchf(pathFormat string, values ...any) error {
	return req.Dof(http.MethodPatch, pathFormat, values...)
}

func (req *Request) DoDelete(path string) error {
	return req.Do(http.MethodDelete, path)
}

func (req *Request) DoDeletef(pathFormat string, values ...any) error {
	return req.Dof(http.MethodDelete, pathFormat, values...)
}

func (req *Request) DoConnect(path string) error {
	return req.Do(http.MethodConnect, path)
}

func (req *Request) DoConnectf(pathFormat string, values ...any) error {
	return req.Dof(http.MethodConnect, pathFormat, values...)
}

func (req *Request) DoOptions(path string) error {
	return req.Do(http.MethodOptions, path)
}

func (req *Request) DoOptionsf(pathFormat string, values ...any) error {
	return req.Dof(http.MethodOptions, pathFormat, values...)
}

func (req *Request) DoTrace(path string) error {
	return req.Do(http.MethodTrace, path)
}

func (req *Request) DoTracef(pathFormat string, values ...any) error {
	return req.Dof(http.MethodTrace, pathFormat, values...)
}

func (req *Request) Dof(method, pathFormat string, values ...any) error {
	path := fmt.Sprintf(pathFormat, values...)
	return req.Do(method, path)
}

func (req *Request) Do(method, path string) error {
	req.Method = method
	req.Path = path
	return req.Client.Overridable.Do(req)
}

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

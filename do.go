package vrest

import (
	"net/http"
)

func Do(req *Request) error {
	err := req.makeHTTPRequest()
	if err != nil {
		return err
	}

	trace := req.Client.traceMaker.New(req)
	defer trace.Close()

	req.Response.Raw, err = req.Overridable.DoHTTPRequest(req)
	defer req.Client.closeRawResponse(req)

	req.Response.Error = req.processHTTPResponse(req.Response.Raw, err)
	trace.OnAfterRequest(req)
	return req.Response.Error
}

func (req *Request) DoGet(path string, pathParams ...string) error {
	return req.Do(http.MethodGet, path, pathParams...)
}

func (req *Request) DoHead(path string, pathParams ...string) error {
	return req.Do(http.MethodHead, path, pathParams...)
}

func (req *Request) DoPost(path string, pathParams ...string) error {
	return req.Do(http.MethodPost, path, pathParams...)
}

func (req *Request) DoPut(path string, pathParams ...string) error {
	return req.Do(http.MethodPut, path, pathParams...)
}

func (req *Request) DoPatch(path string, pathParams ...string) error {
	return req.Do(http.MethodPatch, path, pathParams...)
}

func (req *Request) DoDelete(path string, pathParams ...string) error {
	return req.Do(http.MethodDelete, path, pathParams...)
}

func (req *Request) DoConnect(path string, pathParams ...string) error {
	return req.Do(http.MethodConnect, path, pathParams...)
}

func (req *Request) DoOptions(path string, pathParams ...string) error {
	return req.Do(http.MethodOptions, path, pathParams...)
}

func (req *Request) DoTrace(path string, pathParams ...string) error {
	return req.Do(http.MethodTrace, path, pathParams...)
}

func (req *Request) Do(method, path string, pathParams ...string) error {
	req.Method = method
	req.Path = makePath(path, pathParams...)
	return req.Client.Overridable.Do(req)
}

func DoHTTPRequest(req *Request) (*http.Response, error) {
	return req.Client.httpClient.Do(req.Raw)
}

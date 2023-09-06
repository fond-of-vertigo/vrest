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

	trace := req.Client.traceMaker.NewTrace(req)
	defer trace.End()

	req.Response.Raw, err = req.Overridable.DoHTTPRequest(req)
	defer req.Client.closeRawResponse(req)

	req.Response.Error = req.processHTTPResponse(req.Response.Raw, err)
	trace.OnAfterRequest(req)
	return req.Response.Error
}

func (req *Request) DoGet(pathFormat string, values ...any) error {
	return req.Do(http.MethodGet, pathFormat, values...)
}

func (req *Request) DoHead(pathFormat string, values ...any) error {
	return req.Do(http.MethodHead, pathFormat, values...)
}

func (req *Request) DoPost(pathFormat string, values ...any) error {
	return req.Do(http.MethodPost, pathFormat, values...)
}

func (req *Request) DoPut(pathFormat string, values ...any) error {
	return req.Do(http.MethodPut, pathFormat, values...)
}

func (req *Request) DoPatch(pathFormat string, values ...any) error {
	return req.Do(http.MethodPatch, pathFormat, values...)
}

func (req *Request) DoDelete(pathFormat string, values ...any) error {
	return req.Do(http.MethodDelete, pathFormat, values...)
}

func (req *Request) DoConnect(pathFormat string, values ...any) error {
	return req.Do(http.MethodConnect, pathFormat, values...)
}

func (req *Request) DoOptions(pathFormat string, values ...any) error {
	return req.Do(http.MethodOptions, pathFormat, values...)
}

func (req *Request) DoTrace(pathFormat string, values ...any) error {
	return req.Do(http.MethodTrace, pathFormat, values...)
}

func (req *Request) Do(method, pathFormat string, values ...any) error {
	req.Method = method
	req.Path = fmt.Sprintf(pathFormat, values...)
	return req.Client.Overridable.Do(req)
}

func DoHTTPRequest(req *Request) (*http.Response, error) {
	return req.Client.httpClient.Do(req.Raw)
}

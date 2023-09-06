package vrest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

func Do(req *Request) error {
	err := req.makeHTTPRequest()
	if err != nil {
		return err
	}

	req.Context = context.WithValue(req.Context, "otel.context.ptr", &req.Context)
	req.Client.logger.LogAttrs(req.Context, slog.LevelDebug, "executing http request",
		slog.String("otel.action", "trace.start"),
		slog.String("http.method", req.Raw.Method),
		slog.String("http.url", req.Raw.URL.String()),
		slog.String("http.header", fmt.Sprintf("%v", req.Raw.Header)),
		slog.String("http.body", string(req.BodyBytes)),
	)
	defer req.Client.logger.LogAttrs(req.Context, slog.LevelDebug, "", slog.String("otel.action", "span.end"))

	//trace := req.Client.traceMaker.New(req)
	//defer trace.Close()

	req.Response.Raw, err = req.Overridable.DoHTTPRequest(req)
	defer req.Client.closeRawResponse(req)

	req.Response.Error = req.processHTTPResponse(req.Response.Raw, err)

	//trace.OnAfterRequest(req)
	req.Client.logger.LogAttrs(req.Context, slog.LevelDebug, "executed http request",
		slog.String("otel.action", "span.set_attributes"),
		slog.String("request.error", errorToString(req.Response.Error)),
		slog.String("http.status_code", strconv.Itoa(req.Response.StatusCode())),
		slog.String("http.response_header", fmt.Sprintf("%v", req.Raw.Header)),
		slog.String("http.response_body", string(req.Response.BodyBytes)),
	)

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

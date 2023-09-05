package vrest

import (
	"net/http"
	"reflect"
	"strings"
)

func IsXMLContentType(contentType string) bool {
	return strings.Index(contentType, "/xml") > 0
}

func IsJSONContentType(contentType string) bool {
	return strings.Index(contentType, "/json") > 0
}

func IsSuccess(req *Request) bool {
	statusCode := req.Response.Raw.StatusCode
	return statusCode >= 200 && statusCode < 300
}

func (c *Client) closeRawResponse(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		err := resp.Body.Close()
		if err != nil {
			// TODO: Log here
		}
	}
}

func typeOf(i interface{}) reflect.Type {
	return reflect.Indirect(reflect.ValueOf(i)).Type()
}

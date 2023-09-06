package vrest

import (
	"encoding/base64"
	"log/slog"
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

func (c *Client) closeRawResponse(req *Request) {
	resp := req.Response.Raw
	if resp != nil && resp.Body != nil {
		err := resp.Body.Close()
		if err != nil {
			c.logger.LogAttrs(req.Raw.Context(), slog.LevelError, "error when closing response body", slog.String("error", err.Error()))
		}
	}
}

func typeOf(i interface{}) reflect.Type {
	return reflect.Indirect(reflect.ValueOf(i)).Type()
}

func encodeBasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

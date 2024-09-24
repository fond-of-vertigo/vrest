package vrest

import (
	"encoding/base64"
	"log/slog"
	"reflect"
	"slices"
	"strings"
)

// IsXMLContentType checks whether the content type is XML.
func IsXMLContentType(contentType string) bool {
	return strings.Index(contentType, "/xml") > 0
}

// IsJSONContentType checks whether the content type is JSON.
func IsJSONContentType(contentType string) bool {
	return strings.Index(contentType, "/json") > 0
}

// IsSuccess reports whether the given request status code is a success.
// A success is defined as a status code between 200 and 299.
func IsSuccess(req *Request) bool {
	if req.Response.Raw == nil {
		return false
	}

	statusCode := req.Response.Raw.StatusCode
	if len(req.Response.SuccessStatusCodes) > 0 {
		return slices.Contains(req.Response.SuccessStatusCodes, statusCode)
	}

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

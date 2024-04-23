package vrest

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestRequest_unmarshalResponseBody(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	type args struct {
		urlPath  string
		respBody interface{}
	}
	tests := []struct {
		name        string
		args        args
		want        interface{}
		wantErr     error
		initRequest func(*Request, *args)
	}{{
		name: "unmarshal json time",
		args: args{
			urlPath:  "/unmarshal/json/time",
			respBody: &time.Time{},
		},
		want: toPtr(mustParseTime(t, testTimeValue)),
	}, {
		name: "unmarshal json time failed no content type",
		args: args{
			urlPath:  "/unmarshal/json/time?no-content-type",
			respBody: &time.Time{},
		},
		wantErr: ErrResponseNotUnmarshaled,
	}, {
		name: "unmarshal json time with force content type",
		args: args{
			urlPath:  "/unmarshal/json/time?no-content-type",
			respBody: &time.Time{},
		},
		initRequest: func(req *Request, a *args) {
			req.ForceResponseJSON()
		},
		want: toPtr(mustParseTime(t, testTimeValue)),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewWithClient(ts.Client()).
				SetBaseURL(ts.URL)

			req := c.NewRequest().
				SetResponseBody(tt.args.respBody)
			if tt.initRequest != nil {
				tt.initRequest(req, &tt.args)
			}
			err := req.DoGet(tt.args.urlPath)
			if err != nil {
				if tt.wantErr == nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("unexpected error:\ngot:  %s\nwant: %s", err, tt.wantErr)
				}

				return
			}

			if tt.wantErr != nil {
				t.Fatalf("expected error, got nil")
			}

			if !reflect.DeepEqual(tt.args.respBody, tt.want) {
				t.Fatalf("unexpected response body:\ngot:  %v\nwant: %v", tt.args.respBody, tt.want)
			}
		})
	}
}

func mustParseTime(t *testing.T, value string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		t.Fatalf("failed to parse time: %v", err)
	}
	return parsedTime
}

func toPtr[T any](v T) *T {
	return &v
}

package vrest

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testError struct {
	Message1 string
	Message2 string
}

func (e testError) Error() string {
	return fmt.Sprintf("%s: %s", e.Message1, e.Message2)
}

type testBody struct {
	Text   string
	Number int
}

func TestClient_Do(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"text": "test", "number": 123}`))

		/*w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"message1": "test", "message2": "m2"}`))*/

		/*w.WriteHeader(500)
		_, _ = w.Write([]byte(`Internal Server Error: Something went wrong`))*/
	}))
	defer testServer.Close()

	client := NewWithClient(testServer.Client()).
		SetTraceMaker(&NopTraceMaker{}).
		SetBaseURL(testServer.URL).
		SetErrorBodyType(testError{})

	//client.Overridable.Do = MockDoer(testBody{Text: "text", Number: 456}, nil)
	//client.Overridable.Do = MockDoer(nil, &testError{Message1: "xyz", Message2: "abc"})

	//client.Overridable.DoHTTPRequest = MockHTTPDoer(MockHTTPParams{}, "X-API-Key", "")
	//client.Overridable.DoHTTPRequest = MockHTTPDoer(MockJSONResponse(200, `{"text": "test", "number": 123}`))
	client.Overridable.DoHTTPRequest = MockHTTPDoer(MockJSONResponse(400, `{"message1": "test", "message2": "m2"}`))

	respBody := testBody{}
	err := client.NewRequest().
		SetResponseBody(&respBody).
		DoGet("/test")

	if err != nil {
		var e2 *testError
		if errors.As(err, &e2) {
			_ = fmt.Sprintf("%v", e2)
		}

		t.Fatal(err)
	}
}

/*
type testStruct struct {
	Value string `json:"value"`
}
type errorStruct struct {
	Error string `json:"error"`
}

type testData struct {
	name            string
	call            *Call
	wantErr         bool
	wantReqHeader   map[string]string
	wantRespBody    *testStruct
	wantErrRespBody *errorStruct
}

func TestClient_Execute(t *testing.T) {
	tests := []testData{{
		name: "GET",
		call: &Call{
			Method:  http.MethodPost,
			URL:     "/",
			ReqBody: testStruct{Value: "POST TEST"},
		},
		wantErr:      false,
		wantRespBody: &testStruct{Value: "foobar"},
	}, {
		name: "POST",
		call: &Call{
			Method:  http.MethodPost,
			URL:     "/",
			ReqBody: testStruct{Value: "POST TEST"},
		},
		wantReqHeader: map[string]string{
			"Content-Type": "application/json",
		},
		wantErr:      false,
		wantRespBody: &testStruct{Value: "foobar"},
	}, {
		name: "POST with auth header",
		call: &Call{
			Method:     http.MethodPost,
			URL:        "/",
			AuthHeader: "AuthHeaderValue",
			ReqBody:    testStruct{Value: "POST TEST"},
		},
		wantReqHeader: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "AuthHeaderValue",
		},
		wantErr:      false,
		wantRespBody: &testStruct{Value: "foobar"},
	}, {
		name: "GET with error and empty response body",
		call: &Call{
			Method: http.MethodGet,
			URL:    "/",
		},
		wantErrRespBody: &errorStruct{Error: "Error!"},
		wantErr:         true,
	}, {
		name: "GET with error and error response body",
		call: &Call{
			Method: http.MethodGet,
			URL:    "/",
		},
		wantErr:         true,
		wantErrRespBody: &errorStruct{Error: "Error!"},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given:
			testServer := httptest.NewServer(tt.handleJSONRequest(t))
			defer testServer.Close()

			c := &caller{
				httpClient: testServer.Client(),
			}
			tt.call.URL = testServer.URL + tt.call.URL

			// when:
			tt.call.RespBody = reflect.New(reflect.TypeOf(tt.wantRespBody).Elem()).Interface()
			err := c.Do(tt.call)

			// then:
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantRespBody == nil && tt.call.RespBody == nil {
				// OK
			} else {
				if diff := cmp.Diff(tt.wantRespBody, tt.call.RespBody); diff != "" {
					t.Fatalf(diff)
				}
			}
		})
	}
}

func (tt *testData) handleJSONRequest(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != tt.call.Method {
			t.Fatalf("Unexpected HTTP method '%s'! Expected method '%s'", r.Method, tt.call.Method)
		}

		for k, v := range tt.wantReqHeader {
			if r.Header.Get(k) != v {
				t.Errorf("incorrect value for header '%s'; expected: '%s', actual: '%s'", k, v, r.Header.Get(k))
			}
		}

		if tt.call.ReqBody != nil {
			dec := json.NewDecoder(r.Body)
			decodedReqBody := testStruct{}
			err := dec.Decode(&decodedReqBody)
			if err != nil {
				t.Fatalf("Failed to decode posted request body! Body is: %s", err.Error())
			}

			diff := cmp.Diff(tt.call.ReqBody, decodedReqBody)
			if diff != "" {
				t.Fatalf(diff)
			}
		}

		if tt.wantErr {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			if tt.wantRespBody == nil {
				w.WriteHeader(http.StatusNoContent)
			}
		}
		if tt.wantErrRespBody != nil {
			enc := json.NewEncoder(w)
			err := enc.Encode(&tt.wantErrRespBody)
			if err != nil {
				t.Fatalf("Failed to encode response json: %s", err.Error())
			}
		} else if tt.wantRespBody != nil {
			enc := json.NewEncoder(w)
			err := enc.Encode(&tt.wantRespBody)
			if err != nil {
				t.Fatalf("Failed to encode response json: %s", err.Error())
			}
		}
	}
}
*/

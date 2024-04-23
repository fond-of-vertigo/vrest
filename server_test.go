package vrest

import (
	"net/http"
	"net/http/httptest"
)

const (
	testTimeValue     = `2024-04-23T09:26:44.995288+02:00`
	testJSONTimeValue = `"` + testTimeValue + `"`
)

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /unmarshal/json/time", func(w http.ResponseWriter, r *http.Request) {
		if !r.URL.Query().Has("no-content-type") {
			w.Header().Set("Content-Type", "application/json")
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(testJSONTimeValue))
	})

	return httptest.NewServer(mux)
}

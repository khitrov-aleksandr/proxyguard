package main

import (
	"net/http"
	"net/http/httptest"
)

func startMockServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Is-Mock-Server", "true")
		_, _ = w.Write([]byte("{\"delaySec\":0,\"success\":true}"))
	}))

	return server
}

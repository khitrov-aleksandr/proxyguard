package test_mobile

import (
	"net/http"
	"net/http/httptest"
)

func startMockServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Is-Mock-Server", "true")

		respString := "{\"name\":\"xcom-api\"}"

		if r.RequestURI == "/api/v8/manzana/registration" {
			respString = "{\"token\":\"AAABBBDDD\"}"
		}

		if r.RequestURI == "/api/v8/ecom-auth/login-sms-prestep" {
			respString = "{\"success\":true,\"delaySec\":0}"
		}

		w.Write([]byte(respString))
	}))

	return server
}

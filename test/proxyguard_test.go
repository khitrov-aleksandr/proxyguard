package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Proxyguard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))
	defer server.Close()

	//client := &http.Client{
	//	Transport: &http.Transport{
	//		Proxy: http.ProxyURL(server.URL),
	//	},
	//}

	client := server.Client()

	resp, err := client.Get(server.URL)

	if resp.StatusCode == http.StatusOK {
		t.Fatalf("Bad response code: %d", resp.StatusCode)
	}

	if err != nil {
		t.Fatalf("Failed to get: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
}

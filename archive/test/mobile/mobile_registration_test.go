package test_mobile

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestMobileRegistration(t *testing.T) {
	t.Skip("not ready")
	ctx := context.Background()
	r := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	r.FlushDB(ctx)

	s := startMockServer()
	defer s.Close()

	startProxy(ctx, r, s)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	time.Sleep(1 * time.Second)

	registerTests(t, client)
	//loginTests2(t, client)
}

func registerTests(t *testing.T, c *http.Client) {
	tests := []struct {
		name         string
		email        string
		phone        string
		xDeviceId    string
		xDeviceIdMb  string
		resp         string
		isMockServer string
	}{
		{
			name:         "Pass",
			email:        "qvQ2akbqIO2v@gmail.com",
			phone:        "+79999999999",
			xDeviceId:    "abcd1234",
			xDeviceIdMb:  "efgh5678",
			resp:         "{\"token\":\"AAABBBDDD\"}",
			isMockServer: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			firstReq, _ := http.NewRequest(http.MethodGet, "http://localhost:9999/api/v8", nil)
			firstReq.Header.Add("X-Device-Id", tt.xDeviceId)
			c.Do(firstReq)

			req, _ := http.NewRequest(
				http.MethodPost, "http://localhost:9999/api/v8/manzana/registration",
				strings.NewReader("{\"EmailAddress\":\""+tt.email+"\"}"),
			)

			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("X-Device-Id", tt.xDeviceId)
			req.Header.Add("X-Device-Id-Mb", tt.xDeviceIdMb)

			res, _ := c.Do(req)
			if res.StatusCode != http.StatusOK {
				t.Errorf("Expected %d, got %d", http.StatusOK, res.StatusCode)
			}

			body, _ := io.ReadAll(res.Body)
			defer res.Body.Close()

			if strings.TrimSpace(string(body)) != tt.resp {
				t.Errorf("Expected %s, got %s", tt.resp, string(body))
			}

			if res.Header.Get("Is-Mock-Server") != tt.isMockServer {
				t.Errorf("Expected %s, got %s", tt.isMockServer, res.Header.Get("Is-Mock-Server"))
			}
		})
	}
}

func loginTests2(t *testing.T, c *http.Client) {
	tests := []struct {
		name         string
		phone        string
		xDeviceId    string
		xDeviceIdMb  string
		resp         string
		isMockServer string
	}{
		{
			name:         "Pass first same phone",
			phone:        "79999999999",
			xDeviceId:    "abcd1234",
			xDeviceIdMb:  "efgh5678",
			resp:         "{\"success\":true,\"delaySec\":0}",
			isMockServer: "true",
		},
		{
			name:         "Pass second same phone",
			phone:        "79999999999",
			xDeviceId:    "abcd1234",
			xDeviceIdMb:  "efgh5678",
			resp:         "{\"success\":true,\"delaySec\":0}",
			isMockServer: "true",
		},
		{
			name:         "Block third same phone",
			phone:        "79999999999",
			xDeviceId:    "abcd1234",
			xDeviceIdMb:  "efgh5678",
			resp:         "{\"success\":true,\"delaySec\":0}",
			isMockServer: "",
		},
		{
			name:         "Allow diff phone with same id",
			phone:        "79999999998",
			xDeviceId:    "abcd1234",
			xDeviceIdMb:  "efgh5678",
			resp:         "{\"success\":true,\"delaySec\":0}",
			isMockServer: "true",
		},
		/*{
			name:         "Block by X-Device-Id-Mb header",
			phone:        "79999999997",
			xDeviceId:    "abcd1235",
			xDeviceIdMb:  "",
			resp:         "{\"success\":true,\"delaySec\":0}",
			isMockServer: "",
		},*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			firstReq, _ := http.NewRequest(http.MethodGet, "http://localhost:9999/api/v8", nil)
			firstReq.Header.Add("X-Device-Id", tt.xDeviceId)
			c.Do(firstReq)

			req, _ := http.NewRequest(
				http.MethodPost, "http://localhost:9999/api/v8/ecom-auth/login-sms-prestep",
				strings.NewReader("{\"phone\":\""+tt.phone+"\"}"),
			)

			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("X-Device-Id", tt.xDeviceId)
			req.Header.Add("X-Device-Id-Mb", tt.xDeviceIdMb)

			res, _ := c.Do(req)
			if res.StatusCode != http.StatusOK {
				t.Errorf("Expected %d, got %d", http.StatusOK, res.StatusCode)
			}

			body, _ := io.ReadAll(res.Body)
			defer res.Body.Close()

			if strings.TrimSpace(string(body)) != tt.resp {
				t.Errorf("Expected %s, got %s", tt.resp, string(body))
			}

			if res.Header.Get("Is-Mock-Server") != tt.isMockServer {
				t.Errorf("Expected %s, got %s", tt.isMockServer, res.Header.Get("Is-Mock-Server"))
			}
		})
	}
}

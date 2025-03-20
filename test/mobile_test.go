package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/khitrov-aleksandr/proxyguard/mobile/filter"
	"github.com/khitrov-aleksandr/proxyguard/proxy"
	"github.com/khitrov-aleksandr/proxyguard/repository"
)

/*func init() {
	cfg := config.New()
	rp := repository.NewRedisRepository(redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	}), context.Background())
	go mobile.Run(cfg, rp)
}*/

type RespLoginSmsPrestep struct {
	DelaySec int  `json:"delaySec"`
	Success  bool `json:"success"`
}

func TestLSPRequest(t *testing.T) {
	s := startMockServer()
	defer s.Close()

	startProxy(s)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp := RespLoginSmsPrestep{}

	req, _ := http.NewRequestWithContext(context.Background(),
		http.MethodPost, "http://localhost:9999/api/v8/ecom-auth/login-sms-prestep",
		strings.NewReader("{\"phone\":\"79055902305\"}"),
	)

	req.Header.Add("X-Device-Id-Mb", "1234")
	req.Header.Add("Content-Type", "application/json")

	time.Sleep(1 * time.Second)

	res, err := client.Do(req)

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, res.StatusCode)
	}

	log.Printf("Headers: %v", res.Header)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()

	err = json.Unmarshal(body, &resp)
	if err != nil {
		t.Error(err)
	}

	if resp.DelaySec != 0 {
		t.Errorf("Expected %d, got %d", 0, resp.DelaySec)
	}

	if resp.Success != true {
		t.Errorf("Expected %t, got %t", true, resp.Success)
	}

	t.Logf("Response: %+v", resp)
}

func startMockServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Is-Mock-Server", "true")
		_, _ = w.Write([]byte("{\"delaySec\":0,\"success\":true}"))
	}))

	return server
}

func startProxy(s *httptest.Server) {
	go func() {
		rp := repository.NewRedisRepository(redis.NewClient(&redis.Options{
			Addr: "redis:6379",
		}), context.Background())

		f := filter.New(
			rp,
			logger.NewHandlerLogger("logs/mobile/handle-oz.log"),
		)

		proxy.New(
			"9999",
			s.URL,
			f.Handle,
		).Run()
	}()
}

package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/khitrov-aleksandr/proxyguard/config"
	"github.com/khitrov-aleksandr/proxyguard/mobile"
	"github.com/khitrov-aleksandr/proxyguard/repository"
)

func TestRequest(t *testing.T) {
	cfg := config.New()

	repository := repository.NewRedisRepository(redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	}), context.Background())

	//go site.Run(cfg, repository)
	go mobile.Run(cfg, repository)
	//monitor.Run(cfg)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost, "http://localhost:7082/api/v8/ecom-auth/login-sms-prestep", strings.NewReader("{\"phone\": \"79055902305\"}"))

	if err != nil {
		t.Error(err)
	}

	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, res.StatusCode)
	}

	var data interface{}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		t.Error(err)
	}

	t.Logf("Response: %+v", data)
}

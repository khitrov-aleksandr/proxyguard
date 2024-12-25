package main

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/khitrov-aleksandr/proxyguard/config"
	"github.com/khitrov-aleksandr/proxyguard/mobile"
	"github.com/khitrov-aleksandr/proxyguard/monitor"
	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/khitrov-aleksandr/proxyguard/site"
)

func main() {
	cfg := config.New()

	repository := repository.NewRedisRepository(redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	}), context.Background())

	go site.Run(cfg, repository)
	go mobile.Run(cfg, repository)
	monitor.Run(cfg)
}

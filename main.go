package main

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/khitrov-aleksandr/proxyguard/blocker"
	"github.com/khitrov-aleksandr/proxyguard/config"
	"github.com/khitrov-aleksandr/proxyguard/monitor"
	"github.com/khitrov-aleksandr/proxyguard/proxy"
	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/khitrov-aleksandr/proxyguard/site"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.New()
	e := echo.New()

	repository := repository.NewRedisRepository(redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	}), context.Background())

	blocker := blocker.NewRegisterBlocker(repository)

	go monitor.Run(cfg)
	go site.Run(cfg)

	proxy := proxy.New(cfg, e, blocker)
	proxy.Run()
}

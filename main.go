package main

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/khitrov-aleksandr/proxyguard/blocker"
	"github.com/khitrov-aleksandr/proxyguard/config"
	"github.com/khitrov-aleksandr/proxyguard/proxy"
	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.New()
	e := echo.New()

	repository := repository.NewRedisRepository(redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	}), context.Background())

	blocker := blocker.NewRegisterBlocker(repository)

	proxy := proxy.New(cfg, e, blocker)
	proxy.Run()
}

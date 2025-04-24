package test_mobile

import (
	"context"
	"net/http/httptest"

	"github.com/go-redis/redis/v8"
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/khitrov-aleksandr/proxyguard/mobile/handler"
	"github.com/khitrov-aleksandr/proxyguard/proxy"
	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/labstack/echo/v4"
)

func startProxy(ctx context.Context, r *redis.Client, s *httptest.Server) {
	go func() {
		rp := repository.NewRedisRepository(r, ctx)

		aLog := logger.NewLogger("../../logs/mobile/all-oz.log")
		acLog := logger.NewLogger("../../logs/mobile/accepted-oz.log")

		h := handler.New(
			rp,
			logger.NewHandlerLogger("../../logs/mobile/handle-oz.log"),
		)

		proxy.New(
			"9999",
			s.URL,
			echo.New(),
			h,
			aLog,
			acLog,
		).Run()
	}()
}

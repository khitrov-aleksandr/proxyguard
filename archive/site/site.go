package site

import (
	"github.com/khitrov-aleksandr/proxyguard/config"
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/khitrov-aleksandr/proxyguard/proxy"
	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/khitrov-aleksandr/proxyguard/site/handler"
	"github.com/labstack/echo/v4"
)

func Run(cfg *config.Config, rp repository.Repository) {
	aLog := logger.NewLogger("logs/site/all.log")
	acLog := logger.NewLogger("logs/site/accepted.log")

	h := handler.New(
		rp,
		logger.NewHandlerLogger("logs/site/handle.log"),
	)

	proxy.New(
		cfg.SitePort,
		cfg.SiteBackendUrl,
		echo.New(),
		h,
		aLog,
		acLog,
	).Run()
}

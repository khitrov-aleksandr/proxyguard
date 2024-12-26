package mobile

import (
	"github.com/khitrov-aleksandr/proxyguard/config"
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/khitrov-aleksandr/proxyguard/mobile/handler"
	"github.com/khitrov-aleksandr/proxyguard/proxy"
	"github.com/khitrov-aleksandr/proxyguard/repository"

	"github.com/labstack/echo/v4"
)

func Run(cfg *config.Config, rp repository.Repository) {
	aLog := logger.NewLogger("logs/mobile/all.log")
	acLog := logger.NewLogger("logs/mobile/accepted.log")

	h := handler.New(
		rp,
		logger.NewHandlerLogger("logs/mobile/handle.log"),
	)

	proxy.New(
		cfg.MobilePort,
		cfg.MobileBackendUrl,
		echo.New(),
		h,
		aLog,
		acLog,
	).Run()
}

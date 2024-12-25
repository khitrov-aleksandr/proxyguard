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
	e := echo.New()

	h := handler.New(rp)
	rl := logger.NewLogger("logs/mobile/mobile.log")

	pr := proxy.New(cfg.MobilePort, cfg.MobileBackendUrl, e, h, rl)
	pr.Run()
}

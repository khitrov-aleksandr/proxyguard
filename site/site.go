package site

import (
	"github.com/khitrov-aleksandr/proxyguard/config"
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/khitrov-aleksandr/proxyguard/site/handler"
	"github.com/labstack/echo/v4"
)

func Run(cfg *config.Config, rp repository.Repository) {
	e := echo.New()

	h := handler.New(rp)
	rl := logger.NewLogger("logs/site.log")

	pr := New(cfg.SitePort, cfg.SiteBackendUrl, e, h, rl)
	pr.Run()
}

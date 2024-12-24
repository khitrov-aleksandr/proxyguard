package site

import (
	"fmt"

	"github.com/khitrov-aleksandr/proxyguard/config"
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/labstack/echo/v4"
)

func Run(cfg *config.Config) {
	e := echo.New()

	rl := logger.NewLogger("logs/site.log")

	e.Use(rl.Log)
	e.Start(fmt.Sprintf(":%s", cfg.SitePort))
}

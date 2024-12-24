package monitor

import (
	"fmt"
	"net/http"

	"github.com/khitrov-aleksandr/proxyguard/config"
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/labstack/echo/v4"
)

func Run(cfg *config.Config) {
	e := echo.New()

	rl := logger.NewLogger("logs/monitor.log")

	e.Any("/*", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, GetResponse(), "")
	})

	e.Use(rl.Log)

	e.Start(fmt.Sprintf(":%s", cfg.MonitorPort))
}

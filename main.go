package main

import (
	"github.com/khitrov-aleksandr/proxyguard/config"
	"github.com/khitrov-aleksandr/proxyguard/proxy"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.New()
	e := echo.New()

	proxy := proxy.New(cfg, e)
	proxy.Run()
}

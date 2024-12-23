package proxy

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/khitrov-aleksandr/proxyguard/blocker"
	"github.com/khitrov-aleksandr/proxyguard/config"
	"github.com/khitrov-aleksandr/proxyguard/handler"
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/labstack/echo/v4"
)

type Proxy struct {
	cfg *config.Config
	s   *echo.Echo
	h   *handler.Handler
	rl  *logger.RequestLogger
}

func New(cfg *config.Config, s *echo.Echo, blocker *blocker.RegisterBlocker) *Proxy {
	return &Proxy{
		cfg: cfg,
		s:   s,
		h:   handler.NewHandler(blocker),
		rl:  logger.NewRequestLogger(),
	}
}

func (p *Proxy) Run() {
	url, _ := url.Parse(p.cfg.BackendUrl)
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	p.s.Use(p.h.LogHandler)
	p.s.Use(p.h.RegisterHandler)
	p.s.Use(p.h.LoginHandler)

	p.s.Any("/*", func(c echo.Context) error {
		proxy.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})

	p.s.Start(":" + p.cfg.Port)
}

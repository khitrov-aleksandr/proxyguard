package proxy

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/khitrov-aleksandr/proxyguard/blocker"
	"github.com/khitrov-aleksandr/proxyguard/config"
	"github.com/khitrov-aleksandr/proxyguard/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

type Proxy struct {
	cfg *config.Config
	s   *echo.Echo
	h   *handler.Handler
}

func New(cfg *config.Config, s *echo.Echo, blocker *blocker.RegisterBlocker) *Proxy {
	return &Proxy{
		cfg: cfg,
		s:   s,
		h:   handler.NewHandler(blocker),
	}
}

func (p *Proxy) Run() {
	url, _ := url.Parse(p.cfg.BackendUrl)
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	logFile, _ := os.OpenFile(
		"logs/xcom-sa.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)

	logger := zerolog.New(logFile)

	p.s.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			requestData := make(map[string]interface{})
			_ = json.NewDecoder(c.Request().Body).Decode(&requestData)

			headers, _ := json.Marshal(c.Request().Header)
			body, _ := json.Marshal(requestData)

			logger.Info().
				Timestamp().
				Str("ip", c.RealIP()).
				Str("method", c.Request().Method).
				Str("uri", c.Request().RequestURI).
				RawJSON("headers", headers).
				RawJSON("body", body).
				Int("status", v.Status).
				Msg("")

			return nil
		},
	}))

	p.s.Use(p.h.RegisterHandler)
	p.s.Use(p.h.LoginHandler)

	p.s.Any("/*", func(c echo.Context) error {
		proxy.ServeHTTP(c.Response().Writer, c.Request())

		return nil
	})

	p.s.Start(":" + p.cfg.Port)
}

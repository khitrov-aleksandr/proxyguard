package proxy

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/khitrov-aleksandr/proxyguard/contract"
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/labstack/echo/v4"
)

type Proxy struct {
	port  string
	bUrl  string
	c     *echo.Echo
	h     contract.Handler
	aLog  *logger.Logger
	acLog *logger.Logger
}

func New(
	port string,
	bUrl string,
	c *echo.Echo,
	h contract.Handler,
	aLog *logger.Logger,
	acLog *logger.Logger,
) *Proxy {
	return &Proxy{
		port:  port,
		bUrl:  bUrl,
		c:     c,
		h:     h,
		aLog:  aLog,
		acLog: acLog,
	}
}

func (p *Proxy) Run() {
	p.c.Use(p.aLog.Handler)
	p.c.Use(p.h.Handler)
	p.c.Use(p.acLog.Handler)

	url, _ := url.Parse(p.bUrl)
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	p.c.Any("/*", func(c echo.Context) error {
		proxy.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})

	p.c.Start(":" + p.port)
}

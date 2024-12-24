package site

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/khitrov-aleksandr/proxyguard/contract"
	"github.com/labstack/echo/v4"
)

type Proxy struct {
	port string
	bUrl string
	c    *echo.Echo
	h    contract.Handler
	l    contract.Handler
}

func New(port string, bUrl string, c *echo.Echo, h contract.Handler, l contract.Handler) *Proxy {
	return &Proxy{
		bUrl: bUrl,
		port: port,
		c:    c,
		h:    h,
		l:    l,
	}
}

func (p *Proxy) Run() {
	p.c.Use(p.l.Handler)
	p.c.Use(p.h.Handler)

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

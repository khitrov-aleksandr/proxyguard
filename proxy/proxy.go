package proxy

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi"
)

type Proxy struct {
	r    *chi.Router
	port string
	url  string
}

func New(
	port string,
	url string,
) *Proxy {
	return &Proxy{
		port: port,
		url:  url,
	}
}

func (p *Proxy) Run() {
	url, _ := url.Parse(p.url)
	rProxy := httputil.NewSingleHostReverseProxy(url)

	rProxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	p.r.Any("/*", func(w http.ResponseWriter, r *http.Request) {
		rProxy.ServeHTTP(w, r)
	})

	p.r.Start(":" + p.port)
}

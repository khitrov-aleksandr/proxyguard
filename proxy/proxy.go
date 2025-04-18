package proxy

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi/v5"
)

type Proxy struct {
	port string
	bUrl string
	h    func(http.Handler) http.Handler
}

func NewProxy(
	port string,
	bUrl string,
	h func(http.Handler) http.Handler,
) *Proxy {
	return &Proxy{port, bUrl, h}
}

func (p *Proxy) Run() {
	url, _ := url.Parse(p.bUrl)
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	r := chi.NewRouter()
	r.Use(p.h)
	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	http.ListenAndServe(":"+p.port, r)
}

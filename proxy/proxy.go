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
	//aLog  *logger.Logger
	//acLog *logger.Logger
}

func New(
	port string,
	bUrl string,
	h func(http.Handler) http.Handler,
	//aLog *logger.Logger,
	//acLog *logger.Logger,
) *Proxy {
	return &Proxy{
		port: port,
		bUrl: bUrl,
		h:    h,
		//aLog:  aLog,
		//acLog: acLog,
	}
}

func (p *Proxy) Run() {
	//p.c.Use(p.aLog.Handler)
	//p.c.Use(p.h.Handler)
	//p.c.Use(p.acLog.Handler)

	r := chi.NewRouter()

	r.Use(p.h)

	url, _ := url.Parse(p.bUrl)
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	/*p.c.Any("/*", func(c echo.Context) error {
		proxy.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})

	p.c.Start(":" + p.port)
	*/

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	http.ListenAndServe(":7082", r)
}

package proxy

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"github.com/khitrov-aleksandr/proxyguard/config"
)

func Run(target config.ProxyTarget) {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Recovery())

	reverseProxy := httputil.NewSingleHostReverseProxy(target.Backend)
	reverseProxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	engine.NoRoute(func(c *gin.Context) {
		reverseProxy.ServeHTTP(c.Writer, c.Request)
	})

	addr := ":" + target.Port
	log.Printf("[%s] listening on %s -> %s", target.Name, addr, target.Backend)

	if err := engine.Run(addr); err != nil {
		log.Fatalf("[%s] server failed: %v", target.Name, err)
	}
}

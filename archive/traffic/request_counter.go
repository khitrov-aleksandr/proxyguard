package traffic

import (
	"fmt"

	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/labstack/echo/v4"
)

type GetRequestCounter struct {
	c         echo.Context
	rp        repository.Repository
	method    string
	keyPrefix string
	keySuffix string
	keyExpr   int
}

func NewGetRequestCounter(c echo.Context, rp repository.Repository, method string, keyPrefix string, keyExpr int) *GetRequestCounter {
	return &GetRequestCounter{
		c:         c,
		rp:        rp,
		method:    method,
		keyPrefix: keyPrefix,
		keySuffix: "req",
		keyExpr:   keyExpr,
	}
}

func (g *GetRequestCounter) ByIpAndHeader(hName string) {
	ip := g.c.RealIP()
	req := g.c.Request()

	method := req.Method
	hVal := req.Header.Get(hName)

	g.keySuffix = "ip:header"

	if hVal != "" && method == g.method {
		g.rp.Incr(g.keyName(ip, hVal))
		g.rp.Expr(g.keyName(ip, hVal), g.keyExpr)
	}
}

func (g *GetRequestCounter) keyName(ip string, hVal string) string {
	return fmt.Sprintf("%s:%s:%s:%s", g.keyPrefix, g.keySuffix, ip, hVal)
}

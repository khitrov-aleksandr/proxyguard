package traffic

import (
	"fmt"

	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/labstack/echo/v4"
)

type RequestCounter struct {
	c         echo.Context
	rp        repository.Repository
	method    string
	keyPrefix string
	keyExpr   int
}

func NewRequestCounter(c echo.Context, rp repository.Repository, method string, keyPrefix string, keyExpr int) *RequestCounter {
	return &RequestCounter{
		c:         c,
		rp:        rp,
		method:    method,
		keyPrefix: keyPrefix,
		keyExpr:   keyExpr,
	}
}

func (r *RequestCounter) ByIpAndHeader(hName string) {
	ip := r.c.RealIP()
	req := r.c.Request()

	method := req.Method
	hVal := req.Header.Get(hName)

	key := r.keyName(ip, hVal, "ip:header")

	if hVal != "" && method == r.method {
		r.rp.Incr(key)
		r.rp.Expr(key, r.keyExpr)
	}
}

func (r *RequestCounter) keyName(ip string, hVal string, keySuffix string) string {
	if keySuffix == "" {
		keySuffix = "req"
	}

	return fmt.Sprintf("%s:%s:%s:%s", r.keyPrefix, keySuffix, ip, hVal)
}

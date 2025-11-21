package filter

import (
	"fmt"
	"strconv"

	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/labstack/echo/v4"
)

type RequestFilter struct {
	c         echo.Context
	rp        repository.Repository
	lg        *logger.HandlerLogger
	keyPrefix string
	keySuffix string
}

func NewRequestFilter(c echo.Context, rp repository.Repository, lg *logger.HandlerLogger, keyPrefix string) *RequestFilter {
	return &RequestFilter{
		c:         c,
		rp:        rp,
		lg:        lg,
		keyPrefix: keyPrefix,
		keySuffix: "req",
	}
}

func (r *RequestFilter) ByIpAndHeader(hName string) bool {
	ip := r.c.RealIP()
	req := r.c.Request()

	hVal := req.Header.Get(hName)

	r.keySuffix = "ip:header"

	rCount := r.rp.Get(r.keyName(ip, hVal))
	count, _ := strconv.Atoi(rCount.(string))

	if hVal != "" && count < 2 {
		r.lg.Log(ip, fmt.Sprintf("deny by key: %s count: %d", r.keyName(ip, hVal), count))
		return false
	}

	return true
}

func (r *RequestFilter) keyName(ip string, hVal string) string {
	return fmt.Sprintf("%s:%s:%s:%s", r.keyPrefix, r.keySuffix, ip, hVal)
}

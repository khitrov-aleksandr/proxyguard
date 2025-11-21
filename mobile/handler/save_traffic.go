package handler

import (
	"fmt"
	"strconv"

	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/labstack/echo/v4"
)

func (h *Handler) saveTraffic(c echo.Context, rp repository.Repository) error {
	req := c.Request()
	method := req.Method

	ip := c.RealIP()
	deviceId := req.Header.Get("X-Device-Id")

	if deviceId != "" && method == "GET" {
		getCount := rp.Incr(countKeyName(ip, deviceId))
		rp.Expr(countKeyName(ip, deviceId), 1800)
		h.lg.Log(ip, fmt.Sprintf("save whitelist get method count: %d id: %s", getCount, deviceId))
	}

	return nil
}

func (h *Handler) allowById(c echo.Context, rp repository.Repository) bool {
	req := c.Request()

	ip := c.RealIP()
	deviceId := req.Header.Get("X-Device-Id")

	getCount := rp.Get(countKeyName(ip, deviceId))
	count, _ := strconv.Atoi(getCount.(string))

	if count < 2 {
		h.lg.Log(c.RealIP(), fmt.Sprintf("deny by get method: %s count: %d", deviceId, count))
		return false
	}

	return true
}

func countKeyName(ip string, id string) string {
	return fmt.Sprintf("countRequests:%s:%s", ip, id)
}

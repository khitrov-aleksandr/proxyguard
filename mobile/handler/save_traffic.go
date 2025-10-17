package handler

import (
	"fmt"

	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/labstack/echo/v4"
)

func (h *Handler) saveTraffic(c echo.Context, rp repository.Repository) error {
	req := c.Request()
	method := req.Method

	ip := c.RealIP()
	deviceId := req.Header.Get("X-Device-Id")

	if deviceId != "" && method == "GET" {
		if rp.Get(getWhitelistKey(ip, deviceId)) == "" {
			err := rp.Save(getWhitelistKey(ip, deviceId), deviceId, 1800)
			if err != nil {
				return err
			}

			h.lg.Log(ip, fmt.Sprintf("save whitelist id: %s", deviceId))
		} else {
			rp.Expr(getWhitelistKey(ip, deviceId), 1800)
		}
	}

	return nil
}

func (h *Handler) allowById(c echo.Context, rp repository.Repository) bool {
	req := c.Request()

	ip := c.RealIP()
	deviceId := req.Header.Get("X-Device-Id")

	deviceIdWhitelistKey := rp.Get(getWhitelistKey(ip, deviceId))

	if deviceIdWhitelistKey == "" {
		h.lg.Log(c.RealIP(), fmt.Sprintf("deny by device id: %s whitelist id: %s", deviceId, deviceIdWhitelistKey))
		return false
	}

	return true
}

func getWhitelistKey(ip string, id string) string {
	return fmt.Sprintf("whitelistId:%s:%s", ip, id)
}

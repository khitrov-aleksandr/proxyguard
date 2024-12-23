package handler

import (
	"net/http"

	"github.com/khitrov-aleksandr/proxyguard/faker"
	"github.com/labstack/echo/v4"
)

func (h *Handler) LoginHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		uri := c.Request().RequestURI

		if uri == "/api/v8/ecom-auth/login-sms-prestep" {
			ip := c.RealIP()

			if h.blocker.IsBlocked(ip) {
				h.blocker.Block(ip)
				h.lg.Log(ip, "is blocked")
				return c.JSONPretty(http.StatusOK, faker.GetLoginResponse(), "")
			}

			return next(c)
		}

		return next(c)
	}
}

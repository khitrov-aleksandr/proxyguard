package handler

import (
	"net/http"

	"github.com/khitrov-aleksandr/proxyguard/contract"
	"github.com/khitrov-aleksandr/proxyguard/faker"
	"github.com/khitrov-aleksandr/proxyguard/filter"
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/khitrov-aleksandr/proxyguard/traffic"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	rp repository.Repository
	lg *logger.HandlerLogger
}

func New(rp repository.Repository, lg *logger.HandlerLogger) contract.Handler {
	return &Handler{rp: rp, lg: lg}
}

func (h *Handler) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		traffic.NewMobileTrafficSaver(c, h.rp).Handle()

		uri := c.Request().RequestURI

		if uri == "/api/v8/manzana/registration" || uri == "/mirror/manzana/registration" {
			if !filter.NewMobileFilter(c, h.rp, h.lg).Handle() || h.blockIpByRegister(c, h.rp) {
				return c.JSONPretty(http.StatusOK, faker.GetTokenResponse(), "")
			}
		}

		if uri == "/api/v8/ecom-auth/login-sms-prestep" || uri == "/mirror/ecom-auth/login-sms-prestep" {
			if !filter.NewMobileFilter(c, h.rp, h.lg).Handle() || h.denyLogin(c, h.rp) {
				return c.JSONPretty(http.StatusOK, faker.GetLoginResponse(), "")
			}
		}

		return next(c)
	}
}

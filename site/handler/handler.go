package handler

import (
	"net/http"

	"github.com/khitrov-aleksandr/proxyguard/contract"
	"github.com/khitrov-aleksandr/proxyguard/faker"
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/khitrov-aleksandr/proxyguard/repository"
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
		req := c.Request()
		url := req.RequestURI

		if url == "/api/customer/auth-sms" {
			if h.denySession(c, h.rp) {
				return c.JSONPretty(http.StatusOK, faker.GetAuthSms(), "  ")
			}

			if h.denyCookie(c) {
				return c.JSONPretty(http.StatusOK, faker.GetAuthSms(), "  ")
			}
		}

		return next(c)
	}
}

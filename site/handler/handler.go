package handler

import (
	"net/http"

	"github.com/khitrov-aleksandr/proxyguard/faker"
	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	rp repository.Repository
}

func New(rp repository.Repository) *Handler {
	return &Handler{rp: rp}
}

func (h *Handler) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		url := req.RequestURI

		if url == "/api/customer/auth-sms" {
			if !allowSession(req.Cookies(), req, h.rp) {
				return c.JSONPretty(http.StatusOK, faker.GetAuthSms(), "  ")
			}

			if !allowCookie(req.Cookies()) {
				return c.JSONPretty(http.StatusOK, faker.GetAuthSms(), "  ")
			}
		}

		return next(c)
	}
}

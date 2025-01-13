package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/khitrov-aleksandr/proxyguard/contract"
	"github.com/khitrov-aleksandr/proxyguard/faker"
	"github.com/khitrov-aleksandr/proxyguard/filter"
	"github.com/khitrov-aleksandr/proxyguard/filter/base"
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/labstack/echo/v4"
)

const (
	phoneRusField  string = "phоne"
	phoneField     string = "phone"
	countSmallExpr int64  = 1
	smallExpr      int    = 100
	bigExpr        int    = 86400
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

		prl := filter.NewPhonesRateLimiter(base.NewRateLimiter(h.rp))

		if url == "/api/customer/auth-sms" {
			phone := getPhone(getBody(req))

			if !prl.Allow(phone, countSmallExpr, smallExpr) {
				h.lg.Log(req.RemoteAddr, fmt.Sprintf("block by phone rate limit: phone: %s, count: %d, expr: %d", phone, countSmallExpr, smallExpr))
				return c.JSONPretty(http.StatusOK, faker.GetAuthSms(), "  ")
			}

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

func getPhone(rData map[string]interface{}) string {
	phone := rData[phoneRusField]
	if phone == nil {
		phone = rData[phoneField]
	}

	return phone.(string)
}
func getBody(r *http.Request) (rData map[string]interface{}) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil
	}

	json.Unmarshal(body, &rData)

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	return rData
}

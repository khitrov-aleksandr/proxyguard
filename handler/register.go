package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/khitrov-aleksandr/proxyguard/faker"
	"github.com/khitrov-aleksandr/proxyguard/filter"
	"github.com/labstack/echo/v4"
)

func (h *Handler) RegisterHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		uri := req.RequestURI

		if uri == "/api/v8/manzana/registration" || uri == "/mirror/manzana/registration" {
			requestData := make(map[string]interface{})

			b, _ := io.ReadAll(req.Body)
			json.Unmarshal(b, &requestData)

			req.Body = io.NopCloser(bytes.NewBuffer(b))

			if filter.BlockByEmail(requestData["EmailAddress"].(string)) {
				ip := c.RealIP()
				h.blocker.Block(ip)
				h.lg.Log(ip, "added to black list")
				return c.JSONPretty(http.StatusOK, faker.GetTokenResponse(), "")
			}
		}

		return next(c)
	}
}

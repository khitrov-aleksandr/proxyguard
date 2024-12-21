package handler

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/khitrov-aleksandr/proxyguard/faker"
	"github.com/khitrov-aleksandr/proxyguard/filter"
	"github.com/labstack/echo/v4"
)

func RequestHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		uri := req.RequestURI

		if uri == "/api/v8/manzana/registration" {
			requestData := make(map[string]interface{})

			b, _ := io.ReadAll(req.Body)

			req.Body = io.NopCloser(bytes.NewBuffer(b))

			json.Unmarshal(b, &requestData)

			if filter.BlockByEmail(c) {
				return faker.GetTokenResponse(c)
			}
		}

		return next(c)
	}
}

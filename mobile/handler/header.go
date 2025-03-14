package handler

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) denyHeader(c echo.Context) bool {
	req := c.Request()

	return req.Header.Get("X-Device-Id-Mb") == ""
}

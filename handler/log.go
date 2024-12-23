package handler

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) LogHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		h.rl.Log(c)
		return next(c)
	}
}

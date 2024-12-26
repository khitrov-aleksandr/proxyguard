package handler

import "github.com/labstack/echo/v4"

func (h *Handler) denyCookie(c echo.Context) bool {
	req := c.Request()
	for _, cookie := range req.Cookies() {
		if cookie.Name == "_ym_uid" {
			return false
		}
	}

	h.lg.Log(c.RealIP(), "deny by cookie")
	return true
}

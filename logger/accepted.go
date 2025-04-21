package logger

import "github.com/labstack/echo/v4"

func (l *Logger) AcceptedHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		uri := c.Request().RequestURI

		if uri == "/api/v8/manzana/registration" || uri == "/mirror/manzana/registration" || uri == "/api/v8/ecom-auth/login-sms-prestep" || uri == "/mirror/ecom-auth/login-sms-prestep" {
			l.logWithFormat(c)
		}

		return next(c)
	}
}

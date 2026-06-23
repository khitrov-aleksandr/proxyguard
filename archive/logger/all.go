package logger

import (
	"regexp"

	"github.com/labstack/echo/v4"
)

func (l *Logger) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		uri := c.Request().RequestURI
		matchedForSkipLog, err := regexp.MatchString(`^/mirror*`, uri)
		if err != nil {
			return next(c)
		}

		if !matchedForSkipLog {
			l.logWithFormat(c)
		}

		return next(c)
	}
}

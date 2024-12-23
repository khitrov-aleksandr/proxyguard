package logger

import (
	"encoding/json"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type RequestLogger struct {
	lg zerolog.Logger
}

func NewRequestLogger() *RequestLogger {
	logFile, _ := os.OpenFile(
		"logs/requests.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)

	return &RequestLogger{lg: zerolog.New(logFile)}
}

func (l *RequestLogger) Log(c echo.Context) {
	requestData := make(map[string]interface{})
	_ = json.NewDecoder(c.Request().Body).Decode(&requestData)

	headers, _ := json.Marshal(c.Request().Header)
	body, _ := json.Marshal(requestData)

	l.lg.Info().
		Timestamp().
		Str("ip", c.RealIP()).
		Str("method", c.Request().Method).
		Str("uri", c.Request().RequestURI).
		RawJSON("headers", headers).
		RawJSON("body", body).
		Msg("")
}

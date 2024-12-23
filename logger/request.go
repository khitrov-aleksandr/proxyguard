package logger

import (
	"bytes"
	"encoding/json"
	"io"
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
	req := c.Request()

	requestData := make(map[string]interface{})

	b, _ := io.ReadAll(req.Body)
	json.Unmarshal(b, &requestData)

	req.Body = io.NopCloser(bytes.NewBuffer(b))

	headers, _ := json.Marshal(req.Header)
	body, _ := json.Marshal(requestData)

	l.lg.Info().
		Timestamp().
		Str("ip", c.RealIP()).
		Str("method", req.Method).
		Str("uri", req.RequestURI).
		RawJSON("headers", headers).
		RawJSON("body", body).
		Msg("")
}

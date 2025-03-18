package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Logger struct {
	lg       zerolog.Logger
	filename string
}

func NewLogger(filename string) *Logger {
	logFile, err := os.OpenFile(
		filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)

	if err != nil {
		fmt.Println(err)
	}

	return &Logger{lg: zerolog.New(logFile), filename: filename}
}

func (l *Logger) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		requestData := make(map[string]interface{})

		b, _ := io.ReadAll(req.Body)
		json.Unmarshal(b, &requestData)

		req.Body = io.NopCloser(bytes.NewBuffer(b))

		headers, _ := json.Marshal(req.Header)
		body, _ := json.Marshal(requestData)

		l.lg.Info().
			Timestamp().
			Str("host", req.Host).
			Str("ip", c.RealIP()).
			Str("method", req.Method).
			Str("uri", req.RequestURI).
			RawJSON("headers", headers).
			RawJSON("body", body).
			Msg("")

		return next(c)
	}
}

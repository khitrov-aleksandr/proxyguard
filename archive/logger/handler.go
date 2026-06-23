package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type HandlerLogger struct {
	lg zerolog.Logger
}

func NewHandlerLogger(filename string) *HandlerLogger {
	logFile, _ := os.OpenFile(
		filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)

	return &HandlerLogger{lg: zerolog.New(logFile)}
}

func (l *HandlerLogger) Log(ip string, msg string) {
	l.lg.Info().
		Timestamp().
		Str("ip", ip).
		Msg(msg)
}

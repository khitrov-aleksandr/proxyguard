package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type CustomLogger struct {
	lg zerolog.Logger
}

func NewCustomLogger(filename string) *CustomLogger {
	logFile, _ := os.OpenFile(
		filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)

	return &CustomLogger{lg: zerolog.New(logFile)}
}

func (l *CustomLogger) Log(ip string, msg string) {
	l.lg.Info().
		Timestamp().
		Str("ip", ip).
		Msg(msg)
}

package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type BlockLogger struct {
	lg zerolog.Logger
}

func NewBlockLogger() *BlockLogger {
	logFile, _ := os.OpenFile(
		"logs/block.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)

	return &BlockLogger{lg: zerolog.New(logFile)}
}

func (l *BlockLogger) Log(ip string, msg string) {
	l.lg.Info().
		Timestamp().
		Str("ip", ip).
		Msg(msg)
}

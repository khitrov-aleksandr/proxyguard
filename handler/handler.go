package handler

import (
	"github.com/khitrov-aleksandr/proxyguard/blocker"
	"github.com/khitrov-aleksandr/proxyguard/logger"
)

type Handler struct {
	blocker *blocker.RegisterBlocker
	lg      *logger.BlockLogger
	rl      *logger.RequestLogger
}

func NewHandler(blocker *blocker.RegisterBlocker) *Handler {
	return &Handler{
		blocker: blocker,
		lg:      logger.NewBlockLogger(),
		rl:      logger.NewRequestLogger(),
	}
}

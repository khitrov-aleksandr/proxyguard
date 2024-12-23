package handler

import "github.com/khitrov-aleksandr/proxyguard/blocker"

type Handler struct {
	blocker *blocker.RegisterBlocker
}

func NewHandler(blocker *blocker.RegisterBlocker) *Handler {
	return &Handler{blocker: blocker}
}

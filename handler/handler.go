package handler

import (
	"io"
	"net/http"
)

type handler struct {
}

func New() *handler {
	return &handler{}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

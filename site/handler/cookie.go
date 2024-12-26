package handler

import "net/http"

func (h *Handler) denyCookie(c []*http.Cookie) bool {
	for _, cookie := range c {
		if cookie.Name == "_ym_uid" {
			return false
		}
	}
	return true
}

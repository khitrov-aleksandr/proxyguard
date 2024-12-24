package handler

import "net/http"

func allowCookie(c []*http.Cookie) bool {
	for _, cookie := range c {
		if cookie.Name == "_ym_uid" {
			return true
		}
	}
	return false
}

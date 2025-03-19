package filter

import (
	"net/http"
)

func (f *Filter) denyCookie(r *http.Request) bool {
	for _, cookie := range r.Cookies() {
		if cookie.Name == "_ym_uid" {
			return false
		}
	}

	f.lg.Log(r.RemoteAddr, "deny by cookie")
	return true
}

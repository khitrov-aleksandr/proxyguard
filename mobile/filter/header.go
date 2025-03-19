package filter

import (
	"net/http"
)

func (f *Filter) denyHeader(r *http.Request) bool {
	return r.Header.Get("X-Device-Id-Mb") == ""
}

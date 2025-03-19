package filter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/khitrov-aleksandr/proxyguard/filter"
	"github.com/khitrov-aleksandr/proxyguard/filter/base"
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/khitrov-aleksandr/proxyguard/repository"
)

const (
	phoneRusField  string = "phоne"
	phoneField     string = "phone"
	countSmallExpr int64  = 2
	smallExpr      int    = 100
	bigExpr        int    = 86400
)

type Filter struct {
	rp repository.Repository
	lg *logger.HandlerLogger
}

func New(rp repository.Repository, lg *logger.HandlerLogger) *Filter {
	return &Filter{rp: rp, lg: lg}
}

func (f *Filter) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.RequestURI
		prl := filter.NewPhonesRateLimiter(base.NewRateLimiter(f.rp))

		if url == "/api/customer/auth-sms" {
			phone := getPhone(getBody(r))

			if !prl.Allow(phone, countSmallExpr, smallExpr) {
				f.lg.Log(r.RemoteAddr, fmt.Sprintf("block by phone rate limit: phone: %s, count: %d, expr: %d", phone, countSmallExpr, smallExpr))
				//return c.JSONPretty(http.StatusOK, faker.GetAuthSms(), "  ")
			}

			if f.denySession(r, f.rp) {
				//return c.JSONPretty(http.StatusOK, faker.GetAuthSms(), "  ")
			}

			if f.denyCookie(r) {
				//return c.JSONPretty(http.StatusOK, faker.GetAuthSms(), "  ")
			}
		}

		//return next(c)
	})
}

func getPhone(rData map[string]interface{}) string {
	phone := rData[phoneRusField]
	if phone == nil {
		phone = rData[phoneField]
	}

	return phone.(string)
}
func getBody(r *http.Request) (rData map[string]interface{}) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil
	}

	json.Unmarshal(body, &rData)

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	return rData
}

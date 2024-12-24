package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/khitrov-aleksandr/proxyguard/repository"
)

const (
	sameValCount  = 3
	diffValCount  = 3
	smallExpr     = 120
	bigExpr       = 86400
	phoneRusField = "phоne"
	phoneField    = "phone"
)

func allowSession(c []*http.Cookie, r *http.Request, rp repository.Repository) bool {
	lg := logger.NewCustomLogger("logs/session.log")

	for _, cookie := range c {
		if cookie.Name == "shop_session" {
			session := cookie.Value
			key := getKey(session)
			rData := getBody(r)

			phone := rData[phoneRusField]
			if phone == "" {
				phone = rData[phoneField]
			}

			if rp.Get(key) == "" {
				rp.Save(key, phone, bigExpr)
			} else {
				if rp.Get(key) == phone {
					if rp.Incr(sameValCountSession(session)) > sameValCount {
						lg.Log(r.RemoteAddr, fmt.Sprintf("block as same val: session: %s phone: %s expr %d", session, phone, smallExpr))
						return false
					}

					rp.Expr(sameValCountSession(session), smallExpr)
				} else {
					if rp.Incr(diffValCountSession(session)) > diffValCount {
						lg.Log(r.RemoteAddr, fmt.Sprintf("block as diff val: session: %s phone: %s expr %d", session, phone, bigExpr))
						return false
					}

					rp.Expr(diffValCountSession(session), bigExpr)
				}

				rp.Save(key, phone, bigExpr)
			}
		}
	}

	return true
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

func getKey(session string) string {
	return fmt.Sprintf("shop_session:%s", session)
}

func diffValCountSession(session string) string {
	return fmt.Sprintf("shop_session:diff_val:%s", session)
}

func sameValCountSession(session string) string {
	return fmt.Sprintf("shop_session:same_val:%s", session)
}

package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/labstack/echo/v4"
)

const (
	sameValCount  int64  = 3
	diffValCount  int64  = 3
	smallExpr     int    = 120
	bigExpr       int    = 86400
	phoneRusField string = "phоne"
	phoneField    string = "phone"
)

func (h *Handler) denySession(c echo.Context, rp repository.Repository) bool {
	r := c.Request()
	for _, cookie := range r.Cookies() {
		if cookie.Name == "shop_session" {
			session := cookie.Value
			key := getKey(session)
			rData := getBody(r)

			phone := rData[phoneRusField]
			if phone == nil {
				phone = rData[phoneField]
			}

			if rp.Get(key) == "" {
				rp.Save(key, phone, bigExpr)
			} else {
				if rp.Get(key) == phone {
					if rp.Incr(sameValCountSession(session)) > sameValCount {
						h.lg.Log(r.RemoteAddr, fmt.Sprintf("block as same val: session: %s phone: %s, expr: %d", session, phone, smallExpr))
						return true
					}

					rp.Expr(sameValCountSession(session), smallExpr)
				} else {
					if rp.Incr(diffValCountSession(session)) > diffValCount {
						rp.Incr(sameValCountSession(session))
						rp.Expr(sameValCountSession(session), bigExpr)

						h.lg.Log(r.RemoteAddr, fmt.Sprintf("block as diff val: session: %s phone: %s, expr: %d", session, phone, bigExpr))
						return true
					}

					rp.Expr(diffValCountSession(session), bigExpr)
				}

				rp.Save(key, phone, bigExpr)
			}
		}
	}

	return false
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

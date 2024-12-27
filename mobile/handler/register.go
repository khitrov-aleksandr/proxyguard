package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/khitrov-aleksandr/proxyguard/filter"
	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/labstack/echo/v4"
)

const (
	blockTime      int = 86400
	phoneBlockTime int = 86400
)

func (h *Handler) blockIpByRegister(c echo.Context, rp repository.Repository) bool {
	r := c.Request()
	uri := r.RequestURI

	if uri == "/api/v8/manzana/registration" || uri == "/mirror/manzana/registration" {
		requestData := make(map[string]interface{})

		b, _ := io.ReadAll(r.Body)
		json.Unmarshal(b, &requestData)

		r.Body = io.NopCloser(bytes.NewBuffer(b))

		email := requestData["EmailAddress"].(string)
		phone := requestData["MobilePhone"].(string)

		if filter.BlockByEmail(email) {
			ip := c.RealIP()

			rp.Save(getKey(ip), ip, blockTime)

			rp.Save(getPhoneKey(phone), phone, phoneBlockTime)
			rp.Incr(getPhoneCountKey(phone))
			rp.Expr(getPhoneCountKey(phone), phoneBlockTime)

			h.lg.Log(ip, fmt.Sprintf("block by email: %s", email))
			return true
		}
	}

	return false
}

func (h *Handler) denyLogin(c echo.Context, rp repository.Repository) bool {
	uri := c.Request().RequestURI

	if uri == "/api/v8/ecom-auth/login-sms-prestep" || uri == "/mirror/ecom-auth/login-sms-prestep" {
		ip := c.RealIP()

		if rp.Get(getKey(ip)) == ip {
			rp.Save(getKey(ip), ip, blockTime)

			h.lg.Log(ip, fmt.Sprintf("deny login by ip, expr: %d", blockTime))
			return true
		}
	}

	return false
}

func getKey(ip string) string {
	return fmt.Sprintf("reg_block:%s", ip)
}

func getPhoneKey(phone string) string {
	return fmt.Sprintf("reg_phone:%s", phone)
}

func getPhoneCountKey(phone string) string {
	return fmt.Sprintf("reg_phone:count:%s", phone)
}

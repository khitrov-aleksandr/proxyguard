package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/khitrov-aleksandr/proxyguard/filter"
	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/labstack/echo/v4"
)

const (
	sameValCount       int64 = 2
	diffValCount       int64 = 1
	authCount          int64 = 2
	blockTime          int   = 86400
	phoneBlockTime     int   = 1800
	authPhoneBlockTime int   = 60
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

		phone = strings.ReplaceAll(phone, "+", "")

		if filter.BlockByEmail(email) {
			ip := c.RealIP()

			rp.Save(getKey(ip), ip, blockTime)

			rp.Save(getAuthPhoneKey(phone), 999, authPhoneBlockTime)
			rp.Expr(getAuthPhoneKey(phone), phoneBlockTime)

			h.lg.Log(ip, fmt.Sprintf("block by email: %s", email))
			return true
		}

		if rp.Get(getPhoneKey(phone)) == "" {
			rp.Save(getPhoneKey(phone), email, phoneBlockTime)
		} else {
			if rp.Get(getPhoneKey(phone)) != email {
				if rp.Incr(getPhoneDiffValKey(phone)) > diffValCount {
					rp.Expr(getPhoneDiffValKey(phone), phoneBlockTime)
					h.lg.Log(r.RemoteAddr, fmt.Sprintf("block as diff email by phone: phone: %s, expr: %d", phone, phoneBlockTime))
					return true
				}

				rp.Expr(getPhoneDiffValKey(phone), phoneBlockTime)
			} else {
				if rp.Incr(getPhoneSameValKey(phone)) > sameValCount {
					rp.Expr(getPhoneSameValKey(phone), phoneBlockTime)
					h.lg.Log(r.RemoteAddr, fmt.Sprintf("block as same email by phone: phone: %s, expr: %d", phone, phoneBlockTime))
					return true
				}

				rp.Expr(getPhoneSameValKey(phone), phoneBlockTime)
			}
		}
	}

	return false
}

func (h *Handler) denyLogin(c echo.Context, rp repository.Repository) bool {
	req := c.Request()
	uri := req.RequestURI

	if uri == "/api/v8/ecom-auth/login-sms-prestep" || uri == "/mirror/ecom-auth/login-sms-prestep" {
		ip := c.RealIP()

		if rp.Get(getKey(ip)) == ip {
			rp.Save(getKey(ip), ip, blockTime)

			h.lg.Log(ip, fmt.Sprintf("deny login by ip, expr: %d", blockTime))
			return true
		}

		requestData := make(map[string]interface{})

		b, _ := io.ReadAll(req.Body)
		json.Unmarshal(b, &requestData)

		req.Body = io.NopCloser(bytes.NewBuffer(b))

		phone := requestData["phone"].(string)

		curAuthCount := rp.Incr(getAuthPhoneKey(phone))
		rp.Expr(getAuthPhoneKey(phone), authPhoneBlockTime)

		if curAuthCount > authCount {
			h.lg.Log(ip, fmt.Sprintf("deny login by phone: phone: %s, expr: %d", phone, authPhoneBlockTime))
			return true
		}

		if h.denyHeader(c) {
			h.lg.Log(ip, fmt.Sprintf("deny login by header: phone: %s", phone))
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

func getPhoneSameValKey(phone string) string {
	return fmt.Sprintf("reg_phone:same_val:%s", phone)
}

func getPhoneDiffValKey(phone string) string {
	return fmt.Sprintf("reg_phone:diff_val:%s", phone)
}

func getAuthPhoneKey(phone string) string {
	return fmt.Sprintf("auth_phone:%s", phone)
}

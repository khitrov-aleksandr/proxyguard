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
	blockTime = 86400
)

func blockIpByRegister(c echo.Context, rp repository.Repository) bool {
	r := c.Request()
	uri := r.RequestURI

	if uri == "/api/v8/manzana/registration" || uri == "/mirror/manzana/registration" {
		requestData := make(map[string]interface{})

		b, _ := io.ReadAll(r.Body)
		json.Unmarshal(b, &requestData)

		r.Body = io.NopCloser(bytes.NewBuffer(b))

		if filter.BlockByEmail(requestData["EmailAddress"].(string)) {
			ip := c.RealIP()
			rp.Save(getKey(ip), ip, blockTime)
			return true
		}
	}

	return false
}

func denyLogin(c echo.Context, rp repository.Repository) bool {
	uri := c.Request().RequestURI

	if uri == "/api/v8/ecom-auth/login-sms-prestep" || uri == "/mirror/ecom-auth/login-sms-prestep" {
		ip := c.RealIP()

		if rp.Get(getKey(ip)) == ip {
			rp.Save(getKey(ip), ip, blockTime)
			return true
		}
	}

	return false
}

func getKey(ip string) string {
	return fmt.Sprintf("reg_block:%s", ip)
}

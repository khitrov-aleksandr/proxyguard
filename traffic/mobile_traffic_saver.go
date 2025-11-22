package traffic

import (
	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/labstack/echo/v4"
)

type MobileTrafficSaver struct {
	c  echo.Context
	rp repository.Repository
}

func NewMobileTrafficSaver(c echo.Context, rp repository.Repository) *MobileTrafficSaver {
	return &MobileTrafficSaver{
		c:  c,
		rp: rp,
	}
}

func (m *MobileTrafficSaver) Handle() {
	NewRequestCounter(m.c, m.rp, "GET", "m:get", 1800).ByIpAndHeader("X-Device-Id")
}

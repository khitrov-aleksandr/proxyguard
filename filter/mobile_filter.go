package filter

import (
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/labstack/echo/v4"
)

type MobileFilter struct {
	c  echo.Context
	rp repository.Repository
	lg *logger.HandlerLogger
}

func NewMobileFilter(c echo.Context, rp repository.Repository, lg *logger.HandlerLogger) *MobileFilter {
	return &MobileFilter{
		c:  c,
		rp: rp,
		lg: lg,
	}
}

func (m *MobileFilter) Handle() bool {
	return NewRequestFilter(m.c, m.rp, m.lg, "m:get").ByIpAndHeader("X-Device-Id")
}

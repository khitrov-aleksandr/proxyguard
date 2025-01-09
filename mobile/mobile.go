package mobile

import (
	"github.com/khitrov-aleksandr/proxyguard/config"
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/khitrov-aleksandr/proxyguard/mobile/handler"
	"github.com/khitrov-aleksandr/proxyguard/proxy"
	"github.com/khitrov-aleksandr/proxyguard/repository"

	"github.com/labstack/echo/v4"
)

func Run(cfg *config.Config, rp repository.Repository) {
	go runOz(cfg, rp)
	go runSf(cfg, rp)
	go runSa(cfg, rp)
	go runSt(cfg, rp)
}

func runOz(cfg *config.Config, rp repository.Repository) {
	aLog := logger.NewLogger("logs/mobile/all-oz.log")
	acLog := logger.NewLogger("logs/mobile/accepted-oz.log")

	h := handler.New(
		rp,
		logger.NewHandlerLogger("logs/mobile/handle-oz.log"),
	)

	proxy.New(
		cfg.MobilePortOz,
		cfg.MobileBackendUrlOz,
		echo.New(),
		h,
		aLog,
		acLog,
	).Run()
}

func runSf(cfg *config.Config, rp repository.Repository) {
	aLog := logger.NewLogger("logs/mobile/all-sf.log")
	acLog := logger.NewLogger("logs/mobile/accepted-sf.log")

	h := handler.New(
		rp,
		logger.NewHandlerLogger("logs/mobile/handle-sf.log"),
	)

	proxy.New(
		cfg.MobilePortSf,
		cfg.MobileBackendUrlSf,
		echo.New(),
		h,
		aLog,
		acLog,
	).Run()
}

func runSa(cfg *config.Config, rp repository.Repository) {
	aLog := logger.NewLogger("logs/mobile/all-sa.log")
	acLog := logger.NewLogger("logs/mobile/accepted-sa.log")

	h := handler.New(
		rp,
		logger.NewHandlerLogger("logs/mobile/handle-sa.log"),
	)

	proxy.New(
		cfg.MobilePortSa,
		cfg.MobileBackendUrlSa,
		echo.New(),
		h,
		aLog,
		acLog,
	).Run()
}

func runSt(cfg *config.Config, rp repository.Repository) {
	aLog := logger.NewLogger("logs/mobile/all-st.log")
	acLog := logger.NewLogger("logs/mobile/accepted-st.log")

	h := handler.New(
		rp,
		logger.NewHandlerLogger("logs/mobile/handle-st.log"),
	)

	proxy.New(
		cfg.MobilePortSt,
		cfg.MobileBackendUrlSt,
		echo.New(),
		h,
		aLog,
		acLog,
	).Run()
}

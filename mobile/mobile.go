package mobile

import (
	"github.com/khitrov-aleksandr/proxyguard/config"
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/khitrov-aleksandr/proxyguard/mobile/filter"
	"github.com/khitrov-aleksandr/proxyguard/proxy"
	"github.com/khitrov-aleksandr/proxyguard/repository"
)

func Run(cfg *config.Config, rp repository.Repository) {
	go runOz(cfg, rp)
	go runSf(cfg, rp)
	go runSa(cfg, rp)
	go runSt(cfg, rp)
}

func runOz(cfg *config.Config, rp repository.Repository) {
	f := filter.New(
		rp,
		logger.NewHandlerLogger("logs/mobile/handle-oz.log"),
	)

	proxy.New(
		cfg.MobilePortOz,
		cfg.MobileBackendUrlOz,
		f.Handle,
	).Run()
}

func runSf(cfg *config.Config, rp repository.Repository) {
	f := filter.New(
		rp,
		logger.NewHandlerLogger("logs/mobile/handle-sf.log"),
	)

	proxy.New(
		cfg.MobilePortSf,
		cfg.MobileBackendUrlSf,
		f.Handle,
	).Run()
}

func runSa(cfg *config.Config, rp repository.Repository) {
	f := filter.New(
		rp,
		logger.NewHandlerLogger("../logs/mobile/handle-sa.log"),
	)

	proxy.New(
		cfg.MobilePortSa,
		cfg.MobileBackendUrlSa,
		f.Handle,
	).Run()
}

func runSt(cfg *config.Config, rp repository.Repository) {
	f := filter.New(
		rp,
		logger.NewHandlerLogger("logs/mobile/handle-st.log"),
	)

	proxy.New(
		cfg.MobilePortSt,
		cfg.MobileBackendUrlSt,
		f.Handle,
	).Run()
}

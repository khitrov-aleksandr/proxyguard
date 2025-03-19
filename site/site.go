package site

import (
	"github.com/khitrov-aleksandr/proxyguard/config"
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/khitrov-aleksandr/proxyguard/proxy"
	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/khitrov-aleksandr/proxyguard/site/filter"
)

func Run(cfg *config.Config, rp repository.Repository) {
	//aLog := logger.NewLogger("logs/site/all.log")
	//acLog := logger.NewLogger("logs/site/accepted.log")

	f := filter.New(
		rp,
		logger.NewHandlerLogger("logs/site/handle.log"),
	)

	proxy.New(
		cfg.SitePort,
		cfg.SiteBackendUrl,
		f.Handle,
		//aLog,
		//acLog,
	).Run()
}

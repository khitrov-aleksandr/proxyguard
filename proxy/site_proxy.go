package proxy

import (
	"github.com/khitrov-aleksandr/proxyguard/config"
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/khitrov-aleksandr/proxyguard/site/filter"
)

func RunSite(cfg *config.Config, rp repository.Repository) {
	f := filter.New(
		rp,
		logger.NewHandlerLogger("logs/site/handle.log"),
	)

	NewProxy(
		cfg.SitePort,
		cfg.SiteBackendUrl,
		f.Handle,
	).Run()
}

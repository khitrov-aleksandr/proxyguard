package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/khitrov-aleksandr/proxyguard/config"
	"github.com/khitrov-aleksandr/proxyguard/proxy"
)

func main() {
	cfg := config.New()

	for _, target := range cfg.Site {
		go proxy.Run(target)
	}

	for _, target := range cfg.Mobile {
		go proxy.Run(target)
	}

	log.Printf("started %d site and %d mobile proxy instances", len(cfg.Site), len(cfg.Mobile))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

package config

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var regions = []string{"OZ", "SF", "SA", "ST"}

type ProxyTarget struct {
	Name    string
	Port    string
	Backend *url.URL
}

type Config struct {
	Site        []ProxyTarget
	Mobile      []ProxyTarget
	MonitorPort string
	RedisAddr   string
}

func New() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	cfg := &Config{
		MonitorPort: os.Getenv("MONITOR_PORT"),
		RedisAddr:   os.Getenv("REDIS_ADDR"),
	}

	for _, region := range regions {
		cfg.Site = append(cfg.Site, loadProxy("site", region, "SITE"))
		cfg.Mobile = append(cfg.Mobile, loadProxy("mobile", region, "MOBILE"))
	}

	return cfg
}

func loadProxy(kind, region, prefix string) ProxyTarget {
	port := os.Getenv(fmt.Sprintf("%s_PORT_%s", prefix, region))
	backendURL := os.Getenv(fmt.Sprintf("%s_BACKEND_URL_%s", prefix, region))
	backendPort := os.Getenv(fmt.Sprintf("%s_BACKEND_PORT_%s", prefix, region))

	if port == "" || backendURL == "" {
		log.Fatalf("missing %s_PORT_%s or %s_BACKEND_URL_%s", prefix, region, prefix, region)
	}

	target, err := parseBackendURL(backendURL, backendPort)
	if err != nil {
		log.Fatalf("invalid backend URL for %s_%s: %v", prefix, region, err)
	}

	return ProxyTarget{
		Name:    fmt.Sprintf("%s-%s", kind, strings.ToLower(region)),
		Port:    port,
		Backend: target,
	}
}

func parseBackendURL(backendURL, backendPort string) (*url.URL, error) {
	u, err := url.Parse(backendURL)
	if err != nil {
		return nil, err
	}

	if backendPort != "" {
		u.Host = net.JoinHostPort(u.Hostname(), backendPort)
	}

	return u, nil
}

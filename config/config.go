package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SitePort         string
	SiteBackendUrl   string
	MobilePort       string
	MobileBackendUrl string
	MonitorPort      string
	RedisAddr        string
}

func New() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return &Config{
		SitePort:         os.Getenv("SITE_PORT"),
		SiteBackendUrl:   os.Getenv("SITE_BACKEND_URL"),
		MobilePort:       os.Getenv("MOBILE_PORT"),
		MobileBackendUrl: os.Getenv("MOBILE_BACKEND_URL"),
		MonitorPort:      os.Getenv("MONITOR_PORT"),
		RedisAddr:        os.Getenv("REDIS_ADDR"),
	}
}

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BackendUrl  string
	Port        string
	RedisAddr   string
	SitePort    string
	MonitorPort string
}

func New() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return &Config{
		BackendUrl:  os.Getenv("BACKEND_URL"),
		Port:        os.Getenv("PORT"),
		RedisAddr:   os.Getenv("REDIS_ADDR"),
		SitePort:    os.Getenv("SITE_PORT"),
		MonitorPort: os.Getenv("MONITOR_PORT"),
	}
}

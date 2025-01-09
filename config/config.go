package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SitePort           string
	SiteBackendUrl     string
	MobilePortOz       string
	MobileBackendUrlOz string
	MobilePortSf       string
	MobileBackendUrlSf string
	MobilePortSa       string
	MobileBackendUrlSa string
	MobilePortSt       string
	MobileBackendUrlSt string
	MonitorPort        string
	RedisAddr          string
}

func New() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return &Config{
		SitePort:           os.Getenv("SITE_PORT"),
		SiteBackendUrl:     os.Getenv("SITE_BACKEND_URL"),
		MobilePortOz:       os.Getenv("MOBILE_PORT_OZ"),
		MobileBackendUrlOz: os.Getenv("MOBILE_BACKEND_URL_OZ"),
		MobilePortSf:       os.Getenv("MOBILE_PORT_SF"),
		MobileBackendUrlSf: os.Getenv("MOBILE_BACKEND_URL_SF"),
		MobilePortSa:       os.Getenv("MOBILE_PORT_SA"),
		MobileBackendUrlSa: os.Getenv("MOBILE_BACKEND_URL_SA"),
		MobilePortSt:       os.Getenv("MOBILE_PORT_ST"),
		MobileBackendUrlSt: os.Getenv("MOBILE_BACKEND_URL_ST"),
		MonitorPort:        os.Getenv("MONITOR_PORT"),
		RedisAddr:          os.Getenv("REDIS_ADDR"),
	}
}

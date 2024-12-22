package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BackendUrl string
	Port       string
}

func New() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return &Config{
		BackendUrl: os.Getenv("BACKEND_URL"),
		Port:       os.Getenv("PORT"),
	}
}

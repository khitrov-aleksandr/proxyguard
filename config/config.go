package config

type Config struct {
	Port          string
	BackendScheme string
	BackendHost   string
	BackendPort   string
}

func New() *Config {
	return &Config{
		Port:          "8080",
		BackendScheme: "https",
		BackendHost:   "google.com",
		BackendPort:   "443",
	}
}

package config

type Config struct {
	BackendUrl string
	Port       string
}

func New() *Config {
	return &Config{
		BackendUrl: "https://app-01.prod.superapteka.ru",
		Port:       "8080",
	}
}

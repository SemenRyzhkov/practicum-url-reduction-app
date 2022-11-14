package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Host string `env:"SERVER_ADDRESS"`
}

func New() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

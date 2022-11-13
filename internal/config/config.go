package config

import (
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	//defaultHost         = "localhost:8080"
	defaultWriteTimeout = 5 * time.Second
	defaultReadTimeout  = 5 * time.Second
)

type Config struct {
	Host         string        `env:"SERVER_ADDRESS"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT"`
	//BaseURL      string        `env:"BASE_URL"`
}

func New() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

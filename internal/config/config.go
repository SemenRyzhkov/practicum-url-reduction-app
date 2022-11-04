package config

import (
	"time"
)

const (
	defaultHost         = "localhost:8080"
	defaultWriteTimeout = 5 * time.Second
	defaultReadTimeout  = 5 * time.Second
)

type Config struct {
	Host         string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

func New() Config {
	return Config{
		Host:         defaultHost,
		WriteTimeout: defaultWriteTimeout,
		ReadTimeout:  defaultReadTimeout,
	}
}

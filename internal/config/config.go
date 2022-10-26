package config

import (
	"time"
)

const (
	host         = "localhost:8080"
	writeTimeout = 5 * time.Second
	readTimeout  = 5 * time.Second
)

type Config struct {
	Host         string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

func New() Config {
	return Config{
		Host:         host,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}
}

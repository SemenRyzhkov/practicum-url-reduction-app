package config

import (
	"os"
)

type Config struct {
	Host string
}

func New() (Config, error) {
	return Config{os.Getenv("SERVER_ADDRESS")}, nil
}

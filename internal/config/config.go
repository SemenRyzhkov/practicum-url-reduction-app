package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

type Config struct {
	Host string
}

func New() (Config, error) {
	myEnv, err := godotenv.Read()
	fmt.Println(myEnv)
	if err != nil {
		return Config{}, err
	}
	return Config{Host: myEnv["SERVER_ADDRESS"]}, nil
}

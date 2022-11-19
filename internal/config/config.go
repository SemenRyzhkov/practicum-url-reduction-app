package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

type Config struct {
	Host string `env:"SERVER_ADDRESS"`
}

func New() (Config, error) {
	myEnv, err := godotenv.Read()
	fmt.Println(myEnv)
	if err != nil {
		return Config{}, err
	}
	fmt.Println(myEnv["SERVER_ADDRESS"])
	return Config{Host: myEnv["SERVER_ADDRESS"]}, nil
}

//import (
//	"github.com/caarlos0/env/v6"
//)
//
//type Config struct {
//	Host string `env:"SERVER_ADDRESS"`
//}
//
//func New() (Config, error) {
//	cfg := Config{}
//	err := env.Parse(&cfg)
//	{
//		if err != nil {
//			return Config{}, err
//		}
//	}
//	return cfg, nil
//}

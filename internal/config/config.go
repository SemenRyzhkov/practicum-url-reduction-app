package config

//
//type Config struct {
//	Host string
//}
//
//func New() (Config, error) {
//	myEnv, err := godotenv.Read()
//	fmt.Println(myEnv)
//	if err != nil {
//		return Config{}, err
//	}
//	return Config{Host: myEnv["SERVER_ADDRESS"]}, nil
//}
import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Host string `env:"SERVER_ADDRESS"`
}

func New() (Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	{
		if err != nil {
			return Config{}, err
		}
	}
	return cfg, nil
}

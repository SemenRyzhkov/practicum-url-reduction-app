package config

import (
	"os"
)

//Run SERVER_HOST=$(random domain)
//=== RUN   TestIteration5
//iteration5_test.go:60: Не удалось дождаться пока порт 45195 станет доступен для запроса: context deadline exceeded
//=== RUN   TestIteration5/TestEnvVars
//=== RUN   TestIteration5/TestEnvVars/shorten
//iteration5_test.go:138:
//Error Trace:	/__w/practicum-url-reduction-app/practicum-url-reduction-app/iteration5_test.go:138
///__w/practicum-url-reduction-app/practicum-url-reduction-app/suite.go:91
//Error:      	Received unexpected error:
//Post "http://localhost:45195/": dial tcp 127.0.0.1:45195: connect: connection refused
//Test:       	TestIteration5/TestEnvVars/shorten
//Messages:   	Ошибка при попытке сделать запрос для сокращения URL

//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//если задаю конфиг этим кодом, автотесты валятся со стэктрэйсом выше, хотя
//с точки зрения логики абсолютно ничего не меняетя

type Config struct {
	Host string
}

func New() (Config, error) {
	//myEnv, err := godotenv.Read()
	//if err != nil {
	//	return Config{}, err
	//}
	return Config{Host: os.Getenv("SERVER_ADDRESS")}, nil
}

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

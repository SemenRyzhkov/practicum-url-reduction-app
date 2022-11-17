package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/app"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/common/utils"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/config"
)

//Поддержите конфигурирование сервиса с помощью флагов командной строки наравне с уже имеющимися переменными окружения:
//флаг -a, отвечающий за адрес запуска HTTP-сервера (переменная SERVER_ADDRESS);
//флаг -b, отвечающий за базовый адрес результирующего сокращённого URL (переменная BASE_URL);
//флаг -f, отвечающий за путь до файла с сокращёнными URL (переменная FILE_STORAGE_PATH).

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	utils.HandleFlag()
	flag.Parse()
	fmt.Println(os.Getenv("SERVER_ADDRESS"))
	fmt.Println(os.Getenv("BASE_URL"))
	fmt.Println(os.Getenv("FILE_STORAGE_PATH"))

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	a, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = a.Run()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

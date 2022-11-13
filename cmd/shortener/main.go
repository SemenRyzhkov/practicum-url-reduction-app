package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/app"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/config"
)

//BASE_URL=http://localhost:8080/;SERVER_ADDRESS=localhost:8081

func main() {
	err := godotenv.Load(".env")
	host := os.Getenv("BASE_URL")
	fmt.Println("HOST " + host)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SERVER " + cfg.Host)

	a, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = a.Run()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

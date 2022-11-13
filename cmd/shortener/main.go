package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/app"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/config"
)

func main() {
	os.Setenv("SERVER_ADDRESS", "localhost:8080")
	os.Setenv("BASE_URL", "http://localhost:8080")
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

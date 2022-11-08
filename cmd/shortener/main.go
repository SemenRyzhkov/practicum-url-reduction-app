package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/app"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/config"
)

func main() {
	cfg := config.New()
	a, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = a.Run()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

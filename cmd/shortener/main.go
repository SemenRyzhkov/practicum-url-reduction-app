package main

import (
	"errors"
	"flag"
	"log"
	"net/http"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/app"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/common/utils"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/config"
)

func main() {
	utils.LoadEnvironments(".env")
	utils.HandleFlag()
	flag.Parse()

	serverAddress := utils.GetServerAddress()
	filePath := utils.GetFilePath()
	key := utils.GetKey()
	cfg := config.New(serverAddress, filePath, key)

	a, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = a.Run()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

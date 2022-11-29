package main

import (
	"errors"
	"flag"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/app"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/common/utils"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/config"
)

const (
	host     = "192.168.99.101"
	port     = 5438
	user     = "postgres"
	password = "postgres"
	dbname   = "url_db"
)

func main() {
	utils.LoadEnvironments(".env")

	utils.HandleFlag()
	flag.Parse()

	serverAddress := utils.GetServerAddress()
	dbAddress := utils.GetDBAddress()
	filePath := utils.GetFilePath()
	key := utils.GetKey()
	cfg := config.New(serverAddress, filePath, key, dbAddress)

	a, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = a.Run()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

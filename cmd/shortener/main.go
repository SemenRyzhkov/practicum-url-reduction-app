package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/app"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/common/utils"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/config"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

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

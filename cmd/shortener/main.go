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
	enableHTTPS := utils.GetEnableHTTPS()
	configFilePath := utils.GetConfigFilePath()

	cfg, err := utils.CreateConfig(serverAddress, filePath, key, dbAddress, configFilePath, enableHTTPS)
	if err != nil {
		log.Fatal(err)
	}

	a, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = a.Run(cfg.EnableHTTPS)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

}

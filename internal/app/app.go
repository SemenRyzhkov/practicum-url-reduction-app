package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/common/utils"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/config"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/handlers"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/handlers/dbhandler"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/router"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/cookieservice"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/urlservice"
)

type App struct {
	HTTPServer *http.Server
}

func New(cfg config.Config) (*App, error) {
	log.Println("creating router")
	urlRepository := utils.CreateRepository(cfg.FilePath, cfg.DataBaseAddress)
	urlService := urlservice.New(urlRepository)
	fmt.Println("service success")

	cookieService := cookieservice.New(cfg.Key)
	fmt.Println("cookie service success")
	urlHandler := handlers.NewHandler(urlService, cookieService)
	fmt.Println("cookie  handler success")
	dbHandler := dbhandler.NewDBHandler(cfg.DataBaseAddress)
	fmt.Println("db  handler success")
	urlRouter := router.NewRouter(urlHandler, dbHandler)
	fmt.Println("rou success")

	server := &http.Server{
		Addr:    cfg.Host,
		Handler: urlRouter,
	}

	return &App{
		HTTPServer: server,
	}, nil
}

func (app *App) Run() error {
	log.Println("run server")
	return app.HTTPServer.ListenAndServe()
}

package app

import (
	"log"
	"net/http"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/common/utils"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/config"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/handlers"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/handlers/dbHandler"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/router"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/cookie"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/url"
)

type App struct {
	HTTPServer *http.Server
}

func New(cfg config.Config) (*App, error) {
	log.Println("creating router")
	urlRepository := utils.CreateRepository(cfg.FilePath)
	urlService := url.New(urlRepository)
	cookieService := cookie.New(cfg.Key)
	urlHandler := handlers.NewHandler(urlService, cookieService)
	dbHandler := dbHandler.NewDBHandler(cfg.DataBaseAddress)
	urlRouter := router.NewRouter(urlHandler, dbHandler)

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

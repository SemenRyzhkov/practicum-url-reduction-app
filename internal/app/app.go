package app

import (
	"log"
	"net/http"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/config"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/handlers"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/router"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service"
)

type App struct {
	HTTPServer *http.Server
}

func New(cfg config.Config) (*App, error) {
	log.Println("create router")
	urlRepository := repositories.NewUrlRepository()
	urlService := service.NewUrlService(urlRepository)
	urlHandler := handlers.NewHandler(urlService)
	urlRouter := router.NewRouter(urlHandler)

	server := &http.Server{
		Addr:         cfg.Host,
		Handler:      urlRouter,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
	}

	return &App{
		HTTPServer: server,
	}, nil
}

func (app *App) Run() error {
	log.Println("run server")
	return app.HTTPServer.ListenAndServe()
}

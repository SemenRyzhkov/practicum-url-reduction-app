package app

import (
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/config"
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/router"
	"log"
	"net/http"
)

type App struct {
	HTTPServer *http.Server
}

func New(cfg config.Config) (*App, error) {
	log.Println("create router")
	router := router.NewRouter()

	server := &http.Server{
		Addr:         cfg.Host,
		Handler:      router,
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

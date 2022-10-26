package app

import (
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/config"
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/router"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type App struct {
	HTTPServer *http.Server
}

func New(cfg config.Config) (*App, error) {
	log.Println("create router")
	httpRouter := httprouter.New()
	log.Println("register URL's handler")
	r := router.NewRouter()
	r.Register(httpRouter)

	server := &http.Server{
		Addr:         cfg.Host,
		Handler:      httpRouter,
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

package app

import (
	"log"
	"net/http"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/common/utils"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/config"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/handlers"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/router"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/cookieservice"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/urlservice"
)

type App struct {
	HTTPServer *http.Server
}

func New(cfg config.Config) (*App, error) {
	log.Println("creating router")
	urlRepository, err := utils.CreateRepository(cfg.FilePath, cfg.DataBaseAddress)
	if err != nil {
		return nil, err
	}
	urlService := urlservice.New(urlRepository)
	cookieService, err := cookieservice.New(cfg.Key)
	if err != nil {
		return nil, err
	}
	urlHandler := handlers.NewHandler(urlService, cookieService)
	urlRouter := router.NewRouter(urlHandler)

	server := &http.Server{
		Addr:    cfg.Host,
		Handler: urlRouter,
		//ReadTimeout:  20 * time.Second,
		//WriteTimeout: 20 * time.Second,
	}

	return &App{
		HTTPServer: server,
	}, nil
}

func (app *App) Run() error {
	log.Println("run server")
	return app.HTTPServer.ListenAndServe()
}

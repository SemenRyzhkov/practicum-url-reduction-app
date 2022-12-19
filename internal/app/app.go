package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func New(cfg config.Config, urlService urlservice.URLService) (*App, error) {
	log.Println("creating router")

	cookieService, err := cookieservice.New(cfg.Key)
	if err != nil {
		return nil, err
	}
	urlHandler := handlers.NewHandler(urlService, cookieService)
	urlRouter := router.NewRouter(urlHandler)

	server := &http.Server{
		Addr:         cfg.Host,
		Handler:      urlRouter,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	return &App{
		HTTPServer: server,
	}, nil
}

func (app *App) Close(service urlservice.URLService) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		log.Println("Closeeeeeeeeeee")
		time.Sleep(10 * time.Second)
		app.HTTPServer.Close()
		service.Stop()
	}()

}

func (app *App) Run() error {
	log.Println("run server")
	return app.HTTPServer.ListenAndServe()
}

func CreateService(cfg config.Config) (urlservice.URLService, error) {
	urlRepository, err := utils.CreateRepository(cfg.FilePath, cfg.DataBaseAddress)
	if err != nil {
		return nil, err
	}
	return urlservice.New(urlRepository), nil
}

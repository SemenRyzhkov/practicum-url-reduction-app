package app

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/config"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/handlers"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/fileStorage"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/memoryStorage"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/router"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service"
)

type App struct {
	HTTPServer *http.Server
}

func New(cfg config.Config) (*App, error) {
	log.Println("creating router")
	urlRepository := createRepository()
	urlService := service.NewURLService(urlRepository)
	urlHandler := handlers.NewHandler(urlService)
	urlRouter := router.NewRouter(urlHandler)

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

func createRepository() repositories.URLRepository {
	filePath := os.Getenv("FILE_STORAGE_PATH")
	if len(strings.TrimSpace(filePath)) == 0 {
		return memoryStorage.NewURLMemoryRepository()
	} else {
		return fileStorage.NewURLFileRepository()
	}
}

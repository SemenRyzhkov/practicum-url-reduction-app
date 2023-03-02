package app

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/common/utils"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/config"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/handlers"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/router"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/cookieservice"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/urlservice"
)

// App запускает приложение.
type App struct {
	HTTPServer *http.Server
}

// New конструктор App
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

	//cert, _ := tls.LoadX509KeyPair("localhost.crt", "localhost.key")

	server := &http.Server{
		Addr:         cfg.Host,
		Handler:      urlRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		//TLSConfig: &tls.Config{
		//	Certificates: []tls.Certificate{cert},
		//},
	}
	defer closeHTTPServerAndStopWorkerPool(server, urlRepository)
	return &App{
		HTTPServer: server,
	}, nil
}

func closeHTTPServerAndStopWorkerPool(server *http.Server, repository repositories.URLRepository) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		server.Close()
		repository.StopWorkerPool()
	}()

}

// Run запуск сервера
func (app *App) Run() error {
	log.Println("run server")
	if utils.GetEnableHTTPS() {
		return app.HTTPServer.ListenAndServeTLS("localhost.crt", "localhost.key")
	}
	return app.HTTPServer.ListenAndServe()
}

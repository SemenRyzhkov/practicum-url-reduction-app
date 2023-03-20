package app

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/common/utils"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/config"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/grpcserver"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/urlservice"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/proto"
)

// GRPCApp запускает GRPC приложение.
type GRPCApp struct {
	GRPCServer *grpc.Server
}

// NewGRPC конструктор GRPCApp
func NewGRPC(cfg config.Config) (*GRPCApp, error) {
	log.Println("creating router")
	urlRepository, err := utils.CreateRepository(cfg.FilePath, cfg.DataBaseAddress)
	if err != nil {
		return nil, err
	}
	urlService := urlservice.New(urlRepository)

	urlServerImpl := grpcserver.NewServerImpl(urlService)

	s := grpc.NewServer()

	proto.RegisterURLsServer(s, urlServerImpl)

	return &GRPCApp{
		GRPCServer: s,
	}, nil

}

// Run запуск сервера
func (app *GRPCApp) Run(cfg config.Config) error {
	listen, err := net.Listen("tcp", cfg.Host)
	if err != nil {
		return err
	}
	fmt.Println("Сервер gRPC начал работу")

	// получаем запрос gRPC
	if err := app.GRPCServer.Serve(listen); err != nil {
		return err
	}
	return nil
}

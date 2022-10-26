package service

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type UrlService interface {
	ReduceAndSaveUrl(request *http.Request) (string, error)
	GetUrlById(request *http.Request, params httprouter.Params) (string, error)
}

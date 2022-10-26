package handlers

import (
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var (
	urlService = service.NewUrlService()
)

func ReduceUrl(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	if reduceUrl, err := urlService.ReduceAndSaveUrl(request); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	} else {
		writer.WriteHeader(http.StatusCreated)
		writer.Write([]byte(reduceUrl))
	}
}

func GetUrlById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if url, err := urlService.GetUrlById(request, params); err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
	} else {
		writer.Header().Add("Location", url)
		writer.WriteHeader(http.StatusTemporaryRedirect)
	}
}

package handlers

import (
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/service"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
)

var (
	urlService = service.NewUrlService()
)

func ReduceUrl(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	b, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	if reduceUrl, err := urlService.ReduceAndSaveUrl(string(b[:])); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	} else {
		writer.WriteHeader(http.StatusCreated)
		writer.Write([]byte(reduceUrl))
	}
}

func GetUrlById(writer http.ResponseWriter, _ *http.Request, params httprouter.Params) {
	urlId := params.ByName("id")
	if url, err := urlService.GetUrlById(urlId); err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
	} else {
		writer.Header().Add("Location", url)
		writer.WriteHeader(http.StatusTemporaryRedirect)
	}
}

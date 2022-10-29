package handlers

import (
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/service"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

var (
	urlService = service.NewUrlService()
)

func ReduceUrl(writer http.ResponseWriter, request *http.Request) {
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

func GetUrlById(writer http.ResponseWriter, r *http.Request) {
	urlId := chi.URLParam(r, "id")
	if urlId == "" {
		http.Error(writer, "urlId param is missed", http.StatusBadRequest)
		return
	}
	if url, err := urlService.GetUrlById(urlId); err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
	} else {
		writer.Header().Add("Location", url)
		writer.WriteHeader(http.StatusTemporaryRedirect)
	}
}

package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service"
)

type urlHandlerImpl struct {
	urlService service.UrlService
}

func NewHandler(urlService service.UrlService) UrlHandler {
	return &urlHandlerImpl{urlService}
}

func (h *urlHandlerImpl) ReduceUrl(writer http.ResponseWriter, request *http.Request) {
	b, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	if reduceUrl, err := h.urlService.ReduceAndSaveUrl(string(b)); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	} else {
		writer.WriteHeader(http.StatusCreated)
		writer.Write([]byte(reduceUrl))
	}
}

func (h *urlHandlerImpl) GetUrlById(writer http.ResponseWriter, r *http.Request) {
	urlId := chi.URLParam(r, "id")
	if urlId == "" {
		http.Error(writer, "urlId param is missed", http.StatusBadRequest)
		return
	}
	if url, err := h.urlService.GetUrlById(urlId); err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
	} else {
		writer.Header().Add("Location", url)
		writer.WriteHeader(http.StatusTemporaryRedirect)
	}
}

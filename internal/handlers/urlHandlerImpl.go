package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service"
)

type urlHandlerImpl struct {
	urlService service.URLService
}

func NewHandler(urlService service.URLService) URLHandler {
	return &urlHandlerImpl{urlService}
}

func (h *urlHandlerImpl) ReduceURLTOJSON(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var urlRequest entity.URLRequest
	err := json.NewDecoder(request.Body).Decode(&urlRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	urlResponse, err := h.urlService.ReduceURLToJSON(urlRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	} else {
		writer.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(writer).Encode(urlResponse)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		}
	}
}

func (h *urlHandlerImpl) ReduceURL(writer http.ResponseWriter, request *http.Request) {
	b, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	if reduceURL, err := h.urlService.ReduceAndSaveURL(string(b)); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	} else {
		writer.WriteHeader(http.StatusCreated)
		writer.Write([]byte(reduceURL))
	}
}

func (h *urlHandlerImpl) GetURLByID(writer http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "id")
	if urlID == "" {
		http.Error(writer, "urlID param is missing", http.StatusBadRequest)
		return
	}
	url, err := h.urlService.GetURLByID(urlID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	} else {
		writer.Header().Add("Location", url)
		writer.WriteHeader(http.StatusTemporaryRedirect)
	}
}

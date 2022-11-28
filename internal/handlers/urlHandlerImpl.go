package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/cookie"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/url"
)

type urlHandlerImpl struct {
	urlService    url.URLService
	cookieService cookie.CookieService
}

func (h *urlHandlerImpl) GetAllURL(writer http.ResponseWriter, request *http.Request) {
	userID, cookieErr := h.cookieService.GetUserIDWithCheckCookieAndIssueNewIfCookieIsMissingOrInvalid(writer, request, "userID")
	if cookieErr != nil {
		http.Error(writer, cookieErr.Error(), http.StatusBadRequest)
	}
	userURLList, notFoundErr := h.urlService.GetAllByUserID(userID)
	if notFoundErr != nil {
		http.Error(writer, notFoundErr.Error(), http.StatusNoContent)
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writeErr := json.NewEncoder(writer).Encode(userURLList)
	if writeErr != nil {
		http.Error(writer, writeErr.Error(), http.StatusBadRequest)
	}
}

func NewHandler(urlService url.URLService, cookieService cookie.CookieService) URLHandler {
	return &urlHandlerImpl{urlService, cookieService}
}

func (h *urlHandlerImpl) ReduceURLTOJSON(writer http.ResponseWriter, request *http.Request) {
	userID, cookieErr := h.cookieService.GetUserIDWithCheckCookieAndIssueNewIfCookieIsMissingOrInvalid(writer, request, "userID")
	if cookieErr != nil {
		http.Error(writer, cookieErr.Error(), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	var urlRequest entity.URLRequest
	err := json.NewDecoder(request.Body).Decode(&urlRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	urlResponse, err := h.urlService.ReduceURLToJSON(userID, urlRequest)
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
	userID, cookieErr := h.cookieService.GetUserIDWithCheckCookieAndIssueNewIfCookieIsMissingOrInvalid(writer, request, "userID")
	if cookieErr != nil {
		http.Error(writer, cookieErr.Error(), http.StatusBadRequest)
		return
	}
	b, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	reduceURL, err := h.urlService.ReduceAndSaveURL(userID, string(b))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte(reduceURL))

}

func (h *urlHandlerImpl) GetURLByID(writer http.ResponseWriter, request *http.Request) {
	urlID := chi.URLParam(request, "id")
	if urlID == "" {
		http.Error(writer, "urlID param is missing", http.StatusBadRequest)
		return
	}
	url, err := h.urlService.GetURLByID(urlID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	writer.Header().Add("Location", url)
	writer.WriteHeader(http.StatusTemporaryRedirect)

}

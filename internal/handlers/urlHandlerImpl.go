package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/cookieservice"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/service/urlservice"
)

type urlHandlerImpl struct {
	urlService    urlservice.URLService
	cookieService cookieservice.CookieService
}

func NewHandler(urlService urlservice.URLService, cookieService cookieservice.CookieService) URLHandler {
	return &urlHandlerImpl{urlService, cookieService}
}

func (u *urlHandlerImpl) GetAllURL(writer http.ResponseWriter, request *http.Request) {
	userID, cookieErr := u.cookieService.GetUserIDWithCheckCookieAndIssueNewIfCookieIsMissingOrInvalid(writer, request, "userID")
	if cookieErr != nil {
		http.Error(writer, cookieErr.Error(), http.StatusBadRequest)
		return
	}
	userURLList, notFoundErr := u.urlService.GetAllByUserID(request.Context(), userID)
	if notFoundErr != nil {
		http.Error(writer, notFoundErr.Error(), http.StatusNoContent)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writeErr := json.NewEncoder(writer).Encode(userURLList)
	if writeErr != nil {
		http.Error(writer, writeErr.Error(), http.StatusBadRequest)
		return
	}
}

func (u *urlHandlerImpl) GetURLByID(writer http.ResponseWriter, request *http.Request) {
	urlID := chi.URLParam(request, "id")
	if urlID == "" {
		http.Error(writer, "urlID param is missing", http.StatusBadRequest)
		return
	}
	url, err := u.urlService.GetURLByID(request.Context(), urlID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	writer.Header().Add("Location", url)
	writer.WriteHeader(http.StatusTemporaryRedirect)
}

func (u *urlHandlerImpl) ReduceURLTOJSON(writer http.ResponseWriter, request *http.Request) {
	userID, cookieErr := u.cookieService.GetUserIDWithCheckCookieAndIssueNewIfCookieIsMissingOrInvalid(writer, request, "userID")
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
	urlResponse, err := u.urlService.ReduceURLToJSON(request.Context(), userID, urlRequest)

	if err == nil {
		writer.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(writer).Encode(urlResponse)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}
	writer.WriteHeader(http.StatusConflict)
	err = json.NewEncoder(writer).Encode(urlResponse)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
}

func (u *urlHandlerImpl) ReduceURL(writer http.ResponseWriter, request *http.Request) {
	userID, cookieErr := u.cookieService.GetUserIDWithCheckCookieAndIssueNewIfCookieIsMissingOrInvalid(writer, request, "userID")
	if cookieErr != nil {
		http.Error(writer, cookieErr.Error(), http.StatusBadRequest)
		return
	}
	b, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	reduceURL, err := u.urlService.ReduceAndSaveURL(request.Context(), userID, string(b))
	if err != nil {
		writer.WriteHeader(http.StatusConflict)
		writer.Write([]byte(reduceURL))
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Write([]byte(reduceURL))
}

func (u *urlHandlerImpl) ReduceSeveralURL(writer http.ResponseWriter, request *http.Request) {
	userID, cookieErr := u.cookieService.GetUserIDWithCheckCookieAndIssueNewIfCookieIsMissingOrInvalid(writer, request, "userID")
	if cookieErr != nil {
		http.Error(writer, cookieErr.Error(), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	var urlWithIDRequestList []entity.URLWithIDRequest
	err := json.NewDecoder(request.Body).Decode(&urlWithIDRequestList)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	urlWithIDResponseList, err := u.urlService.ReduceSeveralURL(request.Context(), userID, urlWithIDRequestList)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	} else {
		writer.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(writer).Encode(urlWithIDResponseList)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
	}

}

func (u *urlHandlerImpl) PingConnection(writer http.ResponseWriter, request *http.Request) {
	err := u.urlService.PingConnection()

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

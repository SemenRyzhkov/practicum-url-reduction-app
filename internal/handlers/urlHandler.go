package handlers

import "net/http"

// URLHandler интерфейс для хэндлера
type URLHandler interface {
	GetAllURL(writer http.ResponseWriter, request *http.Request)
	GetURLByID(writer http.ResponseWriter, r *http.Request)

	ReduceURL(writer http.ResponseWriter, request *http.Request)
	ReduceURLTOJSON(writer http.ResponseWriter, request *http.Request)
	ReduceSeveralURL(writer http.ResponseWriter, request *http.Request)

	RemoveAll(writer http.ResponseWriter, request *http.Request)

	PingConnection(writer http.ResponseWriter, request *http.Request)
}

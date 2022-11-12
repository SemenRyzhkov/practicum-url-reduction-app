package handlers

import "net/http"

type URLHandler interface {
	ReduceURL(writer http.ResponseWriter, request *http.Request)
	GetURLByID(writer http.ResponseWriter, r *http.Request)
	ReduceURLTOJSON(writer http.ResponseWriter, request *http.Request)
}

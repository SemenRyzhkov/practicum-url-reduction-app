package handlers

import "net/http"

type UrlHandler interface {
	ReduceUrl(writer http.ResponseWriter, request *http.Request)
	GetUrlById(writer http.ResponseWriter, r *http.Request)
}

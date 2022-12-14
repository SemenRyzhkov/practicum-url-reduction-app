package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/handlers"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/router/middleware"
)

const (
	reduceURLPath        = "/"
	getURLPath           = "/{id}"
	reduceURLToJSONPath  = "/api/shorten"
	allURLPath           = "/api/user/urls"
	pingPath             = "/ping"
	reduceSeveralURLPath = "/api/shorten/batch"
)

func NewRouter(h handlers.URLHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.DecompressRequest, middleware.CompressResponse)
	r.Get(pingPath, h.PingConnection)
	r.Get(getURLPath, h.GetURLByID)
	r.Get(allURLPath, h.GetAllURL)
	r.Post(reduceURLPath, h.ReduceURL)
	r.Post(reduceURLToJSONPath, h.ReduceURLTOJSON)
	r.Post(reduceSeveralURLPath, h.ReduceSeveralURL)
	return r
}

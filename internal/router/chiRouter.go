package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/handlers"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/router/middleware"
)

const (
	reduceURLPath       = "/"
	getURLPath          = "/{id}"
	reduceURLToJSONPath = "/api/shorten"
)

func NewRouter(h handlers.URLHandler) chi.Router {
	r := chi.NewRouter()
	//r.Use(middleware.DecompressRequest)
	r.Use(middleware.GzipHandle)
	r.Get(getURLPath, h.GetURLByID)
	r.Post(reduceURLPath, h.ReduceURL)
	r.Post(reduceURLToJSONPath, h.ReduceURLTOJSON)
	return r
}

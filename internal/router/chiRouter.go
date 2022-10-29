package router

import (
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/handlers"
	"github.com/go-chi/chi/v5"
)

const (
	reduceURLPath = "/"
	getURLPath    = "/{id}"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Get(getURLPath, handlers.GetUrlById)
	r.Post(reduceURLPath, handlers.ReduceUrl)
	return r
}

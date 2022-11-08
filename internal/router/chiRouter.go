package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/handlers"
)

const (
	reduceURLPath = "/"
	getURLPath    = "/{id}"
)

func NewRouter(h handlers.URLHandler) chi.Router {
	r := chi.NewRouter()
	r.Get(getURLPath, h.GetURLByID)
	r.Post(reduceURLPath, h.ReduceURL)
	return r
}

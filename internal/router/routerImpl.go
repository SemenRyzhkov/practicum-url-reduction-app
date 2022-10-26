package router

import (
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/handlers"
	"github.com/julienschmidt/httprouter"
)

var _ Router = &router{}

const (
	reduceURLPath = "/"
	getURLPath    = "/:id"
)

func NewRouter() Router {
	return &router{}
}

type router struct {
}

func (r *router) Register(router *httprouter.Router) {
	router.GET(getURLPath, handlers.GetUrlById)
	router.POST(reduceURLPath, handlers.ReduceUrl)
}

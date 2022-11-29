package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/handlers"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/handlers/dbhandler"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/router/middleware"
)

const (
	reduceURLPath       = "/"
	getURLPath          = "/{id}"
	reduceURLToJSONPath = "/api/shorten"
	allURLPath          = "/api/user/urls"
	pingPath            = "/ping"
)

func NewRouter(h handlers.URLHandler, db dbhandler.DBHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.DecompressRequest, middleware.CompressResponse)
	r.Get(pingPath, db.PingConnection)
	r.Get(getURLPath, h.GetURLByID)
	r.Get(allURLPath, h.GetAllURL)
	r.Post(reduceURLPath, h.ReduceURL)
	r.Post(reduceURLToJSONPath, h.ReduceURLTOJSON)
	return r
}

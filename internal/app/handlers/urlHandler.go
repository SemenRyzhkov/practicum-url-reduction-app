package handlers

import (
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/app/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var _ Handler = &handler{}

const (
	reduceURLPath = "/"
	getURLPath    = "/:id"
)

var (
	urlService = service.NewUrlService()
)

type handler struct {
}

func NewHandler() Handler {
	return &handler{}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(getURLPath, getUrlById)
	router.POST(reduceURLPath, reduceUrl)
}

func reduceUrl(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	if reduceUrl, err := urlService.ReduceAndSaveUrl(request); err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
	} else {
		writer.WriteHeader(http.StatusCreated)
		writer.Write([]byte(reduceUrl))
	}
}

func getUrlById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if url, err := urlService.GetUrlById(request, params); err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
	} else {
		writer.Header().Add("Location", url)
		writer.WriteHeader(307)
	}
}

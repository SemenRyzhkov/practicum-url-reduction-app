package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/handlers"
)

//Добавьте в сервер новый эндпоинт POST /api/shorten, принимающий в теле запроса JSON-объект {"url":"<some_url>"}
//и возвращающий в ответ объект {"result":"<shorten_url>"}.
//Не забудьте добавить тесты на новый эндпоинт, как и на предыдущие.
//Помните про HTTP content negotiation, проставляйте правильные значения в заголовок Content-Type.

const (
	reduceURLPath       = "/"
	getURLPath          = "/{id}"
	reduceURLToJSONPath = "/api/shorten"
)

//Добавьте поддержку gzip в ваш сервис. Научите его:
//принимать запросы в сжатом формате (HTTP-заголовок Content-Encoding);
//отдавать сжатый ответ клиенту, который поддерживает обработку сжатых ответов (HTTP-заголовок Accept-Encoding).
//Вспомните middleware из урока про HTTP-сервер, это может вам помочь.

func NewRouter(h handlers.URLHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(DecompressRequest)
	r.Use(gzipHandle)
	r.Get(getURLPath, h.GetURLByID)
	r.Post(reduceURLPath, h.ReduceURL)
	r.Post(reduceURLToJSONPath, h.ReduceURLTOJSON)
	return r
}

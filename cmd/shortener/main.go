package main

//Напишите сервис для сокращения длинных URL. Требования:

//Сервер должен быть доступен по адресу: http://localhost:8080.

//Сервер должен предоставлять два эндпоинта: POST / и GET /{id}.

//Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и
//сокращённым URL в виде текстовой строки в теле.

//Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ
//с кодом 307 и оригинальным URL в HTTP-заголовке Location.

// Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400.

import (
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/app/handlers"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

const (
	host         = "localhost:8080"
	writeTimeout = 5 * time.Second
	readTimeout  = 5 * time.Second
)

func main() {
	log.Println("create router")
	router := httprouter.New()
	log.Println("register URL's handler")
	handler := handlers.NewHandler()
	handler.Register(router)

	startSever(router)
}

func startSever(router *httprouter.Router) {
	log.Println("start application")

	server := &http.Server{
		Addr:         host,
		Handler:      router,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}

	log.Fatal(server.ListenAndServe())

}

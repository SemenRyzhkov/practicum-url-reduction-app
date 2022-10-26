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
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/app"
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/config"
	"log"
)

func main() {
	cfg := config.New()
	a, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(a.Run())
}

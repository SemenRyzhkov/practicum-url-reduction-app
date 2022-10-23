package service

import (
	"encoding/json"
	"fmt"
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/app/entities"
	"github.com/julienschmidt/httprouter"
	"math/rand"
	"net/http"
)

const localhost = "http://localhost:8080/"

var (
	_          UrlService = &urlServiceImpl{}
	urlStorage            = make(map[string]entities.ReduceUrl)
)

type urlServiceImpl struct {
}

func NewUrlService() UrlService {
	return &urlServiceImpl{}
}

func (u *urlServiceImpl) ReduceAndSaveUrl(request *http.Request) (string, error) {
	var url entities.Url
	json.NewDecoder(request.Body).Decode(&url)

	if isExist(url.Name) {
		return "", fmt.Errorf("url %s already exist", url.Name)
	}

	reduceUrl := mapUrlToReduceUrl(&url)
	saveUrl(&reduceUrl)
	return localhost + reduceUrl.ID, nil
}

func (u *urlServiceImpl) GetUrlById(request *http.Request, params httprouter.Params) (string, error) {
	id := params.ByName("id")

	if url, notFoundErr := findUrlById(id); notFoundErr != nil {
		return "", notFoundErr
	} else {
		return url, nil
	}
}

func findUrlById(id string) (string, error) {
	url, ok := urlStorage[id]
	fmt.Println(url)
	if !ok {
		return "", fmt.Errorf("url with id %d not found", id)
	}
	return url.Name, nil
}

func mapUrlToReduceUrl(url *entities.Url) entities.ReduceUrl {
	id := reducing()
	return entities.ReduceUrl{
		ID:   id,
		Name: url.Name,
	}
}

func saveUrl(reduceUrl *entities.ReduceUrl) {
	fmt.Printf("save url %v\n", reduceUrl)
	urlStorage[reduceUrl.ID] = *reduceUrl
	fmt.Println(urlStorage)
}

func reducing() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func isExist(token string) bool {
	for _, url := range urlStorage {
		if url.Name == token {
			return true
		}
	}
	return false

}

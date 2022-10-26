package service

import (
	"fmt"
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/repositories"
	"github.com/julienschmidt/httprouter"
	"io"
	"math/rand"
	"net/http"
)

const localhost = "http://localhost:8080/"

var (
	_             UrlService = &urlServiceImpl{}
	urlRepository            = repositories.NewUrlRepository()
)

type urlServiceImpl struct {
}

func NewUrlService() UrlService {
	return &urlServiceImpl{}
}

func (u *urlServiceImpl) ReduceAndSaveUrl(request *http.Request) (string, error) {
	b, err := io.ReadAll(request.Body)
	if err != nil {
		return "", fmt.Errorf("not valid request body")
	}

	url := string(b[:])
	reduceUrl := reducing()

	duplicateErr := urlRepository.Save(reduceUrl, url)
	if duplicateErr != nil {
		return "", err
	}
	return localhost + reduceUrl, nil
}

func (u *urlServiceImpl) GetUrlById(request *http.Request, params httprouter.Params) (string, error) {
	id := params.ByName("id")
	return urlRepository.FindById(id)
}

func reducing() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

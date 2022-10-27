package service

import (
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/repositories"
	"math/rand"
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

func (u *urlServiceImpl) ReduceAndSaveUrl(url string) (string, error) {
	reduceUrl := reducing()
	duplicateErr := urlRepository.Save(reduceUrl, url)
	if duplicateErr != nil {
		return "", duplicateErr
	}
	return localhost + reduceUrl, nil
}

func (u *urlServiceImpl) GetUrlById(urlId string) (string, error) {
	return urlRepository.FindById(urlId)
}

func reducing() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

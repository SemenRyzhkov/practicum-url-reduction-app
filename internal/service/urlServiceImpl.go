package service

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
)

const localhost = "http://localhost:8080/"

var (
	_ UrlService = &urlServiceImpl{}
)

type urlServiceImpl struct {
	urlRepository repositories.UrlRepository
}

func NewUrlService(urlRepository repositories.UrlRepository) UrlService {
	return &urlServiceImpl{
		urlRepository,
	}
}

func (u *urlServiceImpl) ReduceAndSaveUrl(url string) (string, error) {
	reduceUrl := reducing(url)
	duplicateErr := u.urlRepository.Save(reduceUrl, url)
	if duplicateErr != nil {
		return "", duplicateErr
	}
	return localhost + reduceUrl, nil
}

func (u *urlServiceImpl) GetUrlById(urlId string) (string, error) {
	return u.urlRepository.FindById(urlId)
}

func reducing(url string) string {
	hash := md5.Sum([]byte(url))
	return hex.EncodeToString(hash[:])
}

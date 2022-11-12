package service

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
)

const localhost = "http://localhost:8080/"

var (
	_ URLService = &urlServiceImpl{}
)

type urlServiceImpl struct {
	urlRepository repositories.URLRepository
}

func NewURLService(urlRepository repositories.URLRepository) URLService {
	return &urlServiceImpl{
		urlRepository,
	}
}

func (u *urlServiceImpl) ReduceUrlToJSON(request entity.URLRequest) (entity.URLResponse, error) {
	reduceURL := reducing(request.URL)
	duplicateErr := u.urlRepository.Save(reduceURL, request.URL)
	if duplicateErr != nil {
		return entity.URLResponse{}, duplicateErr
	}
	return entity.URLResponse{
		Result: localhost + reduceURL,
	}, nil
}

func (u *urlServiceImpl) ReduceAndSaveURL(url string) (string, error) {
	reduceURL := reducing(url)
	duplicateErr := u.urlRepository.Save(reduceURL, url)
	if duplicateErr != nil {
		return "", duplicateErr
	}
	return localhost + reduceURL, nil
}

func (u *urlServiceImpl) GetURLByID(urlID string) (string, error) {
	return u.urlRepository.FindByID(urlID)
}

func reducing(url string) string {
	hash := md5.Sum([]byte(url))
	return hex.EncodeToString(hash[:])
}

package repositories

import (
	"fmt"
	"sync"
)

var _ UrlRepository = &urlRepositoryImpl{}

type urlRepositoryImpl struct {
	mx         sync.Mutex
	urlStorage map[string]string
}

func (u *urlRepositoryImpl) Save(urlId, url string) error {
	u.mx.Lock()
	defer u.mx.Unlock()
	if isExist(u.urlStorage, urlId) {
		return fmt.Errorf("url %s already exist", url)
	}
	u.urlStorage[urlId] = url
	return nil
}

func (u *urlRepositoryImpl) FindById(urlId string) (string, error) {
	u.mx.Lock()
	defer u.mx.Unlock()
	url, ok := u.urlStorage[urlId]
	if !ok {
		return "", fmt.Errorf("url with id %s not found", urlId)
	}
	return url, nil
}

func NewUrlRepository() UrlRepository {
	return &urlRepositoryImpl{
		urlStorage: make(map[string]string),
	}
}

func isExist(urlStorage map[string]string, urlId string) bool {
	if _, ok := urlStorage[urlId]; ok {
		return true
	}
	return false
}

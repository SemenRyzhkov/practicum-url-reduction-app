package repositories

import (
	"fmt"
	"sync"
)

var _ URLRepository = &urlRepositoryImpl{}

type urlRepositoryImpl struct {
	mx         sync.Mutex
	urlStorage map[string]string
}

func (u *urlRepositoryImpl) Save(urlID, url string) error {
	u.mx.Lock()
	defer u.mx.Unlock()
	if isExist(u.urlStorage, urlID) {
		return fmt.Errorf("url %s already exist", url)
	}
	u.urlStorage[urlID] = url
	return nil
}

func (u *urlRepositoryImpl) FindByID(urlID string) (string, error) {
	u.mx.Lock()
	defer u.mx.Unlock()
	url, ok := u.urlStorage[urlID]
	if !ok {
		return "", fmt.Errorf("url with id %s not found", urlID)
	}
	return url, nil
}

func NewURLRepository() URLRepository {
	return &urlRepositoryImpl{
		urlStorage: make(map[string]string),
	}
}

func isExist(urlStorage map[string]string, urlID string) bool {
	if _, ok := urlStorage[urlID]; ok {
		return true
	}
	return false
}

package inmemory

import (
	"fmt"
	"sync"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/mapper"
)

var _ repositories.URLRepository = &urlMemoryRepository{}

type urlMemoryRepository struct {
	mx         sync.Mutex
	urlStorage map[string]map[string]string
}

func (u *urlMemoryRepository) GetAllByUserID(userID string) ([]entity.FullURL, error) {
	userURlMap, ok := u.urlStorage[userID]
	if !ok {
		return nil, fmt.Errorf("user with id %s has not URL's", userID)
	}
	return mapper.FromMapToSliceOfFullURL(userURlMap), nil
}

func (u *urlMemoryRepository) Save(userID, urlID, url string) error {
	u.mx.Lock()
	defer u.mx.Unlock()
	userURLStorage, ok := u.urlStorage[userID]
	if !ok {
		userURLStorage = make(map[string]string)
	}
	if isExist(userURLStorage, urlID) {
		return fmt.Errorf("url %s already exist", url)
	}
	userURLStorage[urlID] = url
	u.urlStorage[userID] = userURLStorage
	return nil
}

func (u *urlMemoryRepository) FindByID(userID, urlID string) (string, error) {
	u.mx.Lock()
	defer u.mx.Unlock()
	userURLStorage, ok := u.urlStorage[userID]
	if !ok {
		return "", fmt.Errorf("user with id %s has not reducing urls", userID)
	}
	url, ok := userURLStorage[urlID]
	if !ok {
		return "", fmt.Errorf("url with id %s not found", urlID)
	}
	return url, nil
}

func New() repositories.URLRepository {
	return &urlMemoryRepository{
		urlStorage: make(map[string]map[string]string),
	}
}

func isExist(urlStorage map[string]string, urlID string) bool {
	_, ok := urlStorage[urlID]
	return ok
}

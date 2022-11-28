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
	mx               sync.Mutex
	commonURLStorage map[string]string
	urlStorage       map[string]map[string]string
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
	u.commonURLStorage[urlID] = url
	return nil
}

func (u *urlMemoryRepository) FindByID(urlID string) (string, error) {
	u.mx.Lock()
	defer u.mx.Unlock()
	url, ok := u.commonURLStorage[urlID]
	if !ok {
		return "", fmt.Errorf("url with id %s not found", urlID)
	}
	return url, nil
}

func New() repositories.URLRepository {
	return &urlMemoryRepository{
		commonURLStorage: make(map[string]string),
		urlStorage:       make(map[string]map[string]string),
	}
}

func isExist(urlStorage map[string]string, urlID string) bool {
	_, ok := urlStorage[urlID]
	return ok
}

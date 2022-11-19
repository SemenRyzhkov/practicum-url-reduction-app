package inmemory

import (
	"fmt"
	"sync"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
)

var _ repositories.URLRepository = &urlMemoryRepository{}

type urlMemoryRepository struct {
	mx         sync.Mutex
	urlStorage map[string]string
}

func (u *urlMemoryRepository) Save(urlID, url string) error {
	u.mx.Lock()
	defer u.mx.Unlock()
	if isExist(u.urlStorage, urlID) {
		return fmt.Errorf("url %s already exist", url)
	}
	u.urlStorage[urlID] = url
	return nil
}

func (u *urlMemoryRepository) FindByID(urlID string) (string, error) {
	u.mx.Lock()
	defer u.mx.Unlock()
	url, ok := u.urlStorage[urlID]
	if !ok {
		return "", fmt.Errorf("url with id %s not found", urlID)
	}
	return url, nil
}

func New() repositories.URLRepository {
	return &urlMemoryRepository{
		urlStorage: make(map[string]string),
	}
}

func isExist(urlStorage map[string]string, urlID string) bool {
	_, ok := urlStorage[urlID]
	return ok
}

package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/urlmapper"
)

var _ repositories.URLRepository = &urlMemoryRepository{}

type urlMemoryRepository struct {
	mx               sync.Mutex
	commonURLStorage map[string]string
	urlStorage       map[string]map[string]string
}

func (u *urlMemoryRepository) RemoveAll(ctx context.Context, removingList []entity.URLDTO) error {
	return nil
}

func (u *urlMemoryRepository) GetAllByUserID(_ context.Context, userID string) ([]entity.FullURL, error) {
	u.mx.Lock()
	userURLMap, ok := u.urlStorage[userID]
	u.mx.Unlock()
	if !ok {
		return nil, fmt.Errorf("user with id %s has not URL's", userID)
	}
	return urlmapper.FromMapToSliceOfFullURL(userURLMap), nil
}

func (u *urlMemoryRepository) Save(_ context.Context, userID, urlID, url string) error {
	u.mx.Lock()
	defer u.mx.Unlock()
	userURLStorage, ok := u.urlStorage[userID]
	if !ok {
		userURLStorage = make(map[string]string)
	}
	if exists(userURLStorage, urlID) {
		return fmt.Errorf("urlservice %s already exist", url)
	}
	userURLStorage[urlID] = url
	u.urlStorage[userID] = userURLStorage
	u.commonURLStorage[urlID] = url
	return nil
}

func (u *urlMemoryRepository) FindByID(_ context.Context, urlID string) (string, error) {
	u.mx.Lock()
	url, ok := u.commonURLStorage[urlID]
	u.mx.Unlock()
	if !ok {
		return "", fmt.Errorf("urlservice with id %s not found", urlID)
	}
	return url, nil
}

func (u *urlMemoryRepository) Ping() error {
	return nil
}

func New() repositories.URLRepository {
	return &urlMemoryRepository{
		commonURLStorage: make(map[string]string),
		urlStorage:       make(map[string]map[string]string),
	}
}

func exists(urlStorage map[string]string, urlID string) bool {
	_, ok := urlStorage[urlID]
	return ok
}

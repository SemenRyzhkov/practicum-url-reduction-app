package inmemory

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity/myerrors"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
)

var (
	_                      repositories.URLRepository = &urlMemoryRepository{}
	ErrRepositoryIsClosing                            = errors.New("repository is closing")
)

type urlKey struct {
	UserID string
	ID     string
}

type urlValue struct {
	OriginalURL string
	Deleted     bool
}

type urlMemoryRepository struct {
	mx         sync.Mutex
	urlStorage map[urlKey]urlValue
	keyList    []urlKey
}

func (u *urlMemoryRepository) StopWorkerPool() {
}

func (u *urlMemoryRepository) RemoveAll(_ context.Context, removingList []entity.URLDTO) error {
	u.mx.Lock()
	u.mx.Unlock()

	for _, dto := range removingList {
		uk := urlKey{
			ID:     dto.ID,
			UserID: dto.UserID,
		}
		u.keyList = append(u.keyList, uk)
	}
	for _, uk := range u.keyList {
		uv := u.urlStorage[uk]
		uv.Deleted = true
		u.urlStorage[uk] = uv
	}

	fmt.Printf("storage after delete %v", u.urlStorage)
	return nil
}

func (u *urlMemoryRepository) GetAllByUserID(_ context.Context, userID string) ([]entity.FullURL, error) {
	listFullURL := make([]entity.FullURL, 0)
	u.mx.Lock()

	for key, value := range u.urlStorage {
		if key.UserID == userID {
			fullURL := entity.FullURL{
				OriginalURL: value.OriginalURL,
				ShortURL:    fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), key.ID),
			}
			listFullURL = append(listFullURL, fullURL)
		}
	}
	u.mx.Unlock()
	if len(listFullURL) == 0 {
		return nil, fmt.Errorf("user with id %s has not URL's", userID)
	}
	return listFullURL, nil
}

func (u *urlMemoryRepository) Save(_ context.Context, userID, urlID, url string) error {
	u.mx.Lock()
	defer u.mx.Unlock()

	uk := urlKey{
		UserID: userID,
		ID:     urlID,
	}

	uv := urlValue{
		OriginalURL: url,
		Deleted:     false,
	}
	if exists(u.urlStorage, uk) {
		return fmt.Errorf("url %s already exists", uv.OriginalURL)
	}
	u.urlStorage[uk] = uv

	return nil
}

func (u *urlMemoryRepository) FindByID(_ context.Context, urlID string) (string, error) {
	var originalURL string
	u.mx.Lock()

	for key, value := range u.urlStorage {
		if key.ID == urlID {
			if value.Deleted {
				deletedErr := myerrors.NewDeletedError(value.OriginalURL, nil)
				return "", deletedErr
			}
			originalURL = value.OriginalURL
		}
	}
	u.mx.Unlock()
	if len(originalURL) == 0 {
		return "", fmt.Errorf("urlservice with id %s not found", urlID)
	}
	return originalURL, nil
}

func (u *urlMemoryRepository) Ping() error {
	return nil
}

func New() repositories.URLRepository {
	return &urlMemoryRepository{
		urlStorage: make(map[urlKey]urlValue),
		keyList:    make([]urlKey, 0),
	}
}

func exists(urlStorage map[urlKey]urlValue, urlKey urlKey) bool {
	_, ok := urlStorage[urlKey]
	return ok
}

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

// variables for urlMemoryRepository
var (
	//URLRepository проверка
	_ repositories.URLRepository = &urlMemoryRepository{}
	//ErrRepositoryIsClosing ошибка закрытия репо
	ErrRepositoryIsClosing = errors.New("repository is closing")
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
	mx         sync.RWMutex
	urlStorage map[urlKey]urlValue
}

// StopWorkerPool остановка воркер-пула
func (u *urlMemoryRepository) StopWorkerPool() {
}

// RemoveAll удаление всех URL
func (u *urlMemoryRepository) RemoveAll(_ context.Context, removingList []entity.URLDTO) error {
	u.mx.Lock()
	defer u.mx.Unlock()

	for _, dto := range removingList {
		uk := urlKey{
			ID:     dto.ID,
			UserID: dto.UserID,
		}
		uv := u.urlStorage[uk]
		uv.Deleted = true
		u.urlStorage[uk] = uv
	}

	return nil
}

// GetAllByUserID поиск всех URL по ID юзера.
func (u *urlMemoryRepository) GetAllByUserID(_ context.Context, userID string) ([]entity.FullURL, error) {
	listFullURL := make([]entity.FullURL, 0)
	u.mx.RLock()

	for key, value := range u.urlStorage {
		if key.UserID == userID {
			fullURL := entity.FullURL{
				OriginalURL: value.OriginalURL,
				ShortURL:    fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), key.ID),
			}
			listFullURL = append(listFullURL, fullURL)
		}
	}
	u.mx.RUnlock()
	if len(listFullURL) == 0 {
		return nil, fmt.Errorf("user with id %s has not URL's", userID)
	}
	return listFullURL, nil
}

// Save сохранение URL.
func (u *urlMemoryRepository) Save(_ context.Context, userID, urlID, url string) error {
	uk := urlKey{
		UserID: userID,
		ID:     urlID,
	}

	uv := urlValue{
		OriginalURL: url,
		Deleted:     false,
	}
	u.mx.Lock()
	if exists(u.urlStorage, uk) {
		u.mx.Unlock()
		return fmt.Errorf("url %s already exists", uv.OriginalURL)
	}
	u.urlStorage[uk] = uv
	u.mx.Unlock()

	return nil
}

// FindByID поиск URL по ID.
func (u *urlMemoryRepository) FindByID(_ context.Context, urlID string) (string, error) {
	var originalURL string
	u.mx.RLock()

	for key, value := range u.urlStorage {
		if key.ID == urlID {
			if value.Deleted {
				u.mx.RUnlock()
				deletedErr := myerrors.NewDeletedError(value.OriginalURL, nil)
				return "", deletedErr
			}
			originalURL = value.OriginalURL
		}
	}
	u.mx.RUnlock()
	if len(originalURL) == 0 {
		return "", fmt.Errorf("url with id %s not found", urlID)
	}
	return originalURL, nil
}

// Ping проверка связи
func (u *urlMemoryRepository) Ping() error {
	return nil
}

// New конструктор.
func New() repositories.URLRepository {
	return &urlMemoryRepository{
		urlStorage: make(map[urlKey]urlValue),
	}
}

// exists существует ли урл
func exists(urlStorage map[urlKey]urlValue, urlKey urlKey) bool {
	_, ok := urlStorage[urlKey]
	return ok
}

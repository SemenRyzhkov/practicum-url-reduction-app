package infile

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/urlmapper"
)

var _ repositories.URLRepository = &urlFileRepository{}

type urlFileRepository struct {
	mx               sync.Mutex
	commonURLStorage map[string]string
	urlStorage       map[string]map[string]string
	consumer         *consumer
	producer         *producer
}

// StopWorkerPool остановка воркер-пула
func (u *urlFileRepository) StopWorkerPool() {
	////TODO implement me
	//panic("implement me")
}

// RemoveAll удаление всех URL
func (u *urlFileRepository) RemoveAll(_ context.Context, _ []entity.URLDTO) error {
	return nil
}

// GetAllByUserID поиск всех URL по ID юзера.
func (u *urlFileRepository) GetAllByUserID(_ context.Context, userID string) ([]entity.FullURL, error) {
	u.mx.Lock()
	userURLMap, ok := u.urlStorage[userID]
	u.mx.Unlock()
	if !ok {
		return nil, fmt.Errorf("user with id %s has not URL's", userID)
	}
	return urlmapper.FromMapToSliceOfFullURL(userURLMap), nil
}

// Save сохранение URL.
func (u *urlFileRepository) Save(_ context.Context, userID, urlID, url string) error {
	u.mx.Lock()
	userURLStorage, ok := u.urlStorage[userID]
	if !ok {
		userURLStorage = make(map[string]string)
	}
	if exists(userURLStorage, urlID) {
		u.mx.Unlock()
		return fmt.Errorf("urlservice %s already exist", url)
	}
	userURLStorage[urlID] = url
	u.urlStorage[userID] = userURLStorage
	u.commonURLStorage[urlID] = url
	u.mx.Unlock()
	savingURL := savingURL{
		UserID: userID,
		URLID:  urlID,
		URL:    url,
	}
	return u.producer.WriteURL(&savingURL)
}

// FindByID поиск URL по ID.
func (u *urlFileRepository) FindByID(_ context.Context, urlID string) (string, error) {
	u.mx.Lock()
	url, ok := u.commonURLStorage[urlID]
	u.mx.Unlock()
	if !ok {
		return "", fmt.Errorf("urlservice with id %s not found", urlID)
	}
	return url, nil
}

// Ping проверка связи
func (u *urlFileRepository) Ping() error {
	return nil
}

// New конструктор.
func New(filePath string) repositories.URLRepository {
	producer, producerErr := NewProducer(filePath)
	if producerErr != nil {
		log.Fatal(producerErr)
	}

	consumer, consumerErr := NewConsumer(filePath)
	if consumerErr != nil {
		log.Fatal(consumerErr)
	}
	defer consumer.Close()

	return &urlFileRepository{
		producer:         producer,
		urlStorage:       consumer.initializeStorage(),
		commonURLStorage: make(map[string]string),
	}
}

func exists(urlStorage map[string]string, urlID string) bool {
	_, ok := urlStorage[urlID]
	return ok
}

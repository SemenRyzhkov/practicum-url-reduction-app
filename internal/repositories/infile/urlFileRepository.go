package infile

import (
	"fmt"
	"log"
	"sync"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories/mapper"
)

var _ repositories.URLRepository = &urlFileRepository{}

type urlFileRepository struct {
	mx               sync.Mutex
	commonURLStorage map[string]string
	urlStorage       map[string]map[string]string
	consumer         *consumer
	producer         *producer
}

func (u *urlFileRepository) GetAllByUserID(userID string) ([]entity.FullURL, error) {
	userURlMap, ok := u.urlStorage[userID]
	if !ok {
		return nil, fmt.Errorf("user with id %s has not URL's", userID)
	}
	return mapper.FromMapToSliceOfFullURL(userURlMap), nil
}

func (u *urlFileRepository) Save(userID, urlID, url string) error {
	u.mx.Lock()
	userURLStorage, ok := u.urlStorage[userID]
	if !ok {
		userURLStorage = make(map[string]string)
	}
	if isExist(userURLStorage, urlID) {
		u.mx.Unlock()
		return fmt.Errorf("url %s already exist", url)
	}
	userURLStorage[urlID] = url
	u.urlStorage[userID] = userURLStorage
	u.commonURLStorage[urlID] = url
	fmt.Println(u.urlStorage)
	u.mx.Unlock()
	savingURL := savingURL{userID, urlID, url}
	return u.producer.WriteURL(&savingURL)
}

func (u *urlFileRepository) FindByID(urlID string) (string, error) {
	u.mx.Lock()
	url, ok := u.commonURLStorage[urlID]
	u.mx.Unlock()
	if !ok {
		return "", fmt.Errorf("url with id %s not found", urlID)
	}
	return url, nil
}

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

func isExist(urlStorage map[string]string, urlID string) bool {
	_, ok := urlStorage[urlID]
	return ok
}

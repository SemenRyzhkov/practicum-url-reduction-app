package infile

import (
	"fmt"
	"log"
	"sync"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
)

var _ repositories.URLRepository = &urlFileRepository{}

type urlFileRepository struct {
	mx         sync.Mutex
	urlStorage map[string]string
	consumer   *consumer
	producer   *producer
}

func (u *urlFileRepository) Save(urlID, url string) error {
	u.mx.Lock()
	if isExist(u.urlStorage, urlID) {
		u.mx.Unlock()
		return fmt.Errorf("url %s already exist", url)
	}
	u.urlStorage[urlID] = url
	u.mx.Unlock()
	savingURL := savingURL{urlID, url}
	return u.producer.WriteURL(&savingURL)
}

func (u *urlFileRepository) FindByID(urlID string) (string, error) {
	u.mx.Lock()
	url, ok := u.urlStorage[urlID]
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
		producer:   producer,
		urlStorage: consumer.initializeStorage(),
	}
}

func isExist(urlStorage map[string]string, urlID string) bool {
	_, ok := urlStorage[urlID]
	return ok
}

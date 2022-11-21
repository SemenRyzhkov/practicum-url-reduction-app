package infile

import (
	"fmt"
	"log"
	"os"
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
	savingURL := savingURL{urlID, url}
	u.mx.Unlock()
	return u.producer.WriteURL(&savingURL)
}

func (u *urlFileRepository) FindByID(urlID string) (string, error) {
	u.mx.Lock()
	url, ok := u.urlStorage[urlID]
	if !ok {
		u.mx.Unlock()
		return "", fmt.Errorf("url with id %s not found", urlID)
	}
	u.mx.Unlock()
	return url, nil
}

func New() repositories.URLRepository {
	filePath := os.Getenv("FILE_STORAGE_PATH")

	producer, producerErr := NewProducer(filePath)
	if producerErr != nil {
		log.Fatal(producerErr)
	}

	consumer, consumerErr := NewConsumer(filePath)
	if consumerErr != nil {
		log.Fatal(consumerErr)
	}

	return &urlFileRepository{
		producer:   producer,
		urlStorage: consumer.initializeStorage(),
	}
}

func isExist(urlStorage map[string]string, urlID string) bool {
	_, ok := urlStorage[urlID]
	return ok
}

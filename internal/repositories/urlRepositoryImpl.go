package repositories

import "fmt"

var (
	_          UrlRepository = &urlRepositoryImpl{}
	urlStorage               = make(map[string]string)
)

type urlRepositoryImpl struct {
}

func (u urlRepositoryImpl) Save(urlId, url string) error {
	if isExist(url) {
		return fmt.Errorf("url %s already exist", url)
	}
	urlStorage[urlId] = url
	return nil
}

func (u urlRepositoryImpl) FindById(urlId string) (string, error) {
	url, ok := urlStorage[urlId]
	if !ok {
		return "", fmt.Errorf("url with id %d not found", urlId)
	}
	return url, nil
}

func NewUrlRepository() UrlRepository {
	return &urlRepositoryImpl{}
}

func isExist(token string) bool {
	for _, url := range urlStorage {
		if url == token {
			return true
		}
	}
	return false
}

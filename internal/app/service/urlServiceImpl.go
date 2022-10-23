package service

import (
	"encoding/json"
	"fmt"
	"github.com/SemenRyzhkov/practicum-url-reduction-app.git/internal/app/entities"
	"github.com/julienschmidt/httprouter"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

var (
	_          UrlService = &urlServiceImpl{}
	idCounter  uint64
	urlStorage = make(map[uint64]entities.ReduceUrl)
)

type urlServiceImpl struct {
}

func NewUrlService() UrlService {
	return &urlServiceImpl{}
}

func (u *urlServiceImpl) ReduceAndSaveUrl(request *http.Request) (string, error) {
	var url entities.Url
	json.NewDecoder(request.Body).Decode(&url)

	if !isValidUrl(url.Name) {
		return "", fmt.Errorf("wrong url %s", url.Name)
	}

	if isExist(url.Name) {
		return "", fmt.Errorf("url %s already exist", url.Name)
	}

	reduceUrl := mapUrlToReduceUrl(&url)
	saveUrl(reduceUrl)
	return reduceUrl.ReduceName, nil
}

func (u *urlServiceImpl) GetUrlById(request *http.Request, params httprouter.Params) (string, error) {
	if id, parsingErr := getIdFromParams(params); parsingErr != nil {
		return "", parsingErr
	} else {
		if url, notFoundErr := findUrlById(id); notFoundErr != nil {
			return "", notFoundErr
		} else {
			return url, nil
		}
	}
}

func findUrlById(id int) (string, error) {
	url, ok := urlStorage[uint64(id)]
	fmt.Println(url)
	if !ok {
		return "", fmt.Errorf("url with id %d not found", id)
	}
	return url.Name, nil
}

func getIdFromParams(params httprouter.Params) (int, error) {
	idVal := params.ByName("id")
	id, err := strconv.Atoi(idVal)
	if err != nil {
		return 0, fmt.Errorf("parameter id not integer")
	}
	return id, nil
}

func mapUrlToReduceUrl(url *entities.Url) entities.ReduceUrl {
	reduceName := reducing()
	idCounter++
	return entities.ReduceUrl{
		ID:         idCounter,
		Name:       url.Name,
		ReduceName: reduceName,
	}
}

func saveUrl(reduceUrl entities.ReduceUrl) {
	fmt.Printf("save url %v", reduceUrl)
	urlStorage[reduceUrl.ID] = reduceUrl
	fmt.Println(urlStorage)
}

func reducing() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func isValidUrl(token string) bool {
	_, err := url.ParseRequestURI(token)
	if err != nil {
		return false
	}
	u, err := url.Parse(token)
	if err != nil || u.Host == "" {
		return false
	}
	return true
}

func isExist(token string) bool {
	for _, url := range urlStorage {
		if url.Name == token {
			return true
		}
	}
	return false
}

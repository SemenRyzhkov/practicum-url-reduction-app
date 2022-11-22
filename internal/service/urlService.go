package service

import (
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
)

type URLService interface {
	ReduceAndSaveURL(url string) (string, error)
	GetURLByID(urlID string) (string, error)
	ReduceURLToJSON(request entity.URLRequest) (entity.URLResponse, error)
}

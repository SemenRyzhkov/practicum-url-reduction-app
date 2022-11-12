package service

import (
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
)

type URLService interface {
	ReduceAndSaveURL(url string) (string, error)
	GetURLByID(urlID string) (string, error)
	ReduceUrlToJSON(request entity.URLRequest) entity.URLResponse
}

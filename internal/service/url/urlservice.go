package url

import (
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
)

type URLService interface {
	ReduceAndSaveURL(userID, url string) (string, error)
	GetURLByID(urlID string) (string, error)
	ReduceURLToJSON(userID string, request entity.URLRequest) (entity.URLResponse, error)
	GetAllByUserID(userID string) ([]entity.FullURL, error)
}

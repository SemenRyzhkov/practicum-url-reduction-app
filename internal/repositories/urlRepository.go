package repositories

import (
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
)

type URLRepository interface {
	Save(userID, urlID, url string) error
	FindByID(userID, urlID string) (string, error)
	GetAllByUserID(userID string) ([]entity.FullURL, error)
}

package repositories

import (
	"context"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
)

type URLRepository interface {
	Save(ctx context.Context, userID, urlID, url string) error
	FindByID(ctx context.Context, urlID string) (string, error)
	GetAllByUserID(ctx context.Context, userID string) ([]entity.FullURL, error)
	RemoveAll(ctx context.Context, removingListChannel []entity.URLDTO) error
	StopWorkerPool()
	Ping() error
}

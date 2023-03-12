package repositories

import (
	"context"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
)

// URLRepository интерфейс для репозитория
type URLRepository interface {
	Save(ctx context.Context, userID, urlID, url string) error
	FindByID(ctx context.Context, urlID string) (string, error)
	GetAllByUserID(ctx context.Context, userID string) ([]entity.FullURL, error)
	GetStats(ctx context.Context) (entity.Stats, error)
	RemoveAll(ctx context.Context, removingListChannel []entity.URLDTO) error
	StopWorkerPool()
	Ping() error
}

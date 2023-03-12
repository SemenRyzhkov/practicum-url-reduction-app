package urlservice

import (
	"context"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
)

// URLService интерфейс для урл-сервиса
type URLService interface {
	GetURLByID(ctx context.Context, urlID string) (string, error)
	GetAllByUserID(ctx context.Context, userID string) ([]entity.FullURL, error)
	ReduceAndSaveURL(ctx context.Context, userID, url string) (string, error)
	ReduceURLToJSON(ctx context.Context, userID string, request entity.URLRequest) (entity.URLResponse, error)
	ReduceSeveralURL(ctx context.Context, userID string, list []entity.URLWithIDRequest) ([]entity.URLWithIDResponse, error)
	RemoveAll(ctx context.Context, userID string, removingList []string) error
	PingConnection() error
	GetStats(ctx context.Context) (entity.Stats, error)
}

package urlservice

import (
	"context"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
)

type URLService interface {
	ReduceAndSaveURL(ctx context.Context, userID, url string) (string, error)
	GetURLByID(ctx context.Context, urlID string) (string, error)
	ReduceURLToJSON(ctx context.Context, userID string, request entity.URLRequest) (entity.URLResponse, error)
	GetAllByUserID(ctx context.Context, userID string) ([]entity.FullURL, error)
}

package urlservice

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
)

var _ URLService = &urlServiceImpl{}

type urlServiceImpl struct {
	urlRepository repositories.URLRepository
}

func New(urlRepository repositories.URLRepository) URLService {
	return &urlServiceImpl{
		urlRepository,
	}
}

func (u *urlServiceImpl) GetAllByUserID(ctx context.Context, userID string) ([]entity.FullURL, error) {
	return u.urlRepository.GetAllByUserID(ctx, userID)
}

func (u *urlServiceImpl) ReduceURLToJSON(ctx context.Context, userID string, request entity.URLRequest) (entity.URLResponse, error) {
	reduceURL := reducing(request.URL)
	duplicateErr := u.urlRepository.Save(ctx, userID, reduceURL, request.URL)
	if duplicateErr != nil {
		return entity.URLResponse{}, duplicateErr
	}
	return entity.URLResponse{Result: fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), reduceURL)}, nil
}

func (u *urlServiceImpl) ReduceAndSaveURL(ctx context.Context, userID, url string) (string, error) {
	reduceURL := reducing(url)
	duplicateErr := u.urlRepository.Save(ctx, userID, reduceURL, url)
	if duplicateErr != nil {
		return "", duplicateErr
	}
	return fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), reduceURL), nil
}

func (u *urlServiceImpl) GetURLByID(ctx context.Context, urlID string) (string, error) {
	return u.urlRepository.FindByID(ctx, urlID)
}

func reducing(url string) string {
	hash := md5.Sum([]byte(url))
	return hex.EncodeToString(hash[:])
}

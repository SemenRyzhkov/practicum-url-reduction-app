package urlservice

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
)

var _ URLService = &urlServiceImpl{}

type urlServiceImpl struct {
	urlRepository repositories.URLRepository
}

// New конструктор.
func New(urlRepository repositories.URLRepository) URLService {
	return &urlServiceImpl{
		urlRepository,
	}
}

// GetAllByUserID поиск всех URL по ID юзера.
func (u *urlServiceImpl) GetAllByUserID(ctx context.Context, userID string) ([]entity.FullURL, error) {
	return u.urlRepository.GetAllByUserID(ctx, userID)
}

// GetURLByID поиск URL по ID.
func (u *urlServiceImpl) GetURLByID(ctx context.Context, urlID string) (string, error) {
	return u.urlRepository.FindByID(ctx, urlID)
}

// ReduceURLToJSON сокращение и сохранение URL.
func (u *urlServiceImpl) ReduceURLToJSON(ctx context.Context, userID string, request entity.URLRequest) (entity.URLResponse, error) {
	reduceURL := Reducing(request.URL)
	duplicateErr := u.urlRepository.Save(ctx, userID, reduceURL, request.URL)
	if duplicateErr != nil {
		return entity.URLResponse{}, duplicateErr
	}
	return entity.URLResponse{Result: fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), reduceURL)}, nil
}

// ReduceAndSaveURL сокращение и сохранение URL.
func (u *urlServiceImpl) ReduceAndSaveURL(ctx context.Context, userID, url string) (string, error) {
	reduceURL := Reducing(url)
	duplicateErr := u.urlRepository.Save(ctx, userID, reduceURL, url)
	if duplicateErr != nil {
		return "", duplicateErr
	}
	return fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), reduceURL), nil
}

// ReduceSeveralURL сокращение и сохранение нескольких URL.
func (u *urlServiceImpl) ReduceSeveralURL(ctx context.Context, userID string, list []entity.URLWithIDRequest) ([]entity.URLWithIDResponse, error) {
	var urlWithIDResponseList []entity.URLWithIDResponse
	//TODO by statement solution
	for _, urlReq := range list {
		correlationID := urlReq.CorrelationID
		reduceURL := Reducing(urlReq.OriginalURL)
		duplicateErr := u.urlRepository.Save(ctx, userID, reduceURL, urlReq.OriginalURL)
		if duplicateErr != nil {
			return nil, duplicateErr
		}
		urlWihIDResponse := entity.URLWithIDResponse{
			CorrelationID: correlationID,
			ShortURL:      fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), reduceURL),
		}
		urlWithIDResponseList = append(urlWithIDResponseList, urlWihIDResponse)
	}
	return urlWithIDResponseList, nil
}

// RemoveAll удаление всех URL
func (u *urlServiceImpl) RemoveAll(ctx context.Context, userID string, removingList []string) error {
	now := time.Now()
	defer func() {
		fmt.Println(time.Since(now))
	}()

	removingDTOList := make([]entity.URLDTO, 0)

	for _, u := range removingList {
		ud := entity.URLDTO{
			ID:      u,
			UserID:  userID,
			Deleted: true,
		}
		removingDTOList = append(removingDTOList, ud)
	}
	return u.urlRepository.RemoveAll(ctx, removingDTOList)
}

func (u *urlServiceImpl) GetStats(ctx context.Context) (entity.Stats, error) {
	return u.urlRepository.GetStats(ctx)
}

// PingConnection проверка связи
func (u *urlServiceImpl) PingConnection() error {
	return u.urlRepository.Ping()
}

// Reducing самый главный метод во всем приложении
func Reducing(url string) string {
	hash := md5.Sum([]byte(url))
	return hex.EncodeToString(hash[:])
}

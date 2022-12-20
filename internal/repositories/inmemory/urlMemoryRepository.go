package inmemory

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity/myerrors"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
)

var (
	_                      repositories.URLRepository = &urlMemoryRepository{}
	ErrRepositoryIsClosing                            = errors.New("repository is closing")
)

type urlMemoryRepository struct {
	mx            sync.RWMutex
	urlStorage    []entity.URLDTO
	deletionQueue chan entity.URLDTO
	done          chan struct{}
	wg            sync.WaitGroup
	once          sync.Once
}

func (u *urlMemoryRepository) StopWorkerPool() {
	//u.once.Do(func() {
	//	close(u.done)
	//})
	//
	//u.once.Do(func() {
	//	close(u.deletionQueue)
	//})
	//
	//u.wg.Wait()
}

func (u *urlMemoryRepository) addURLToDeletionQueue(ud entity.URLDTO) error {
	select {
	case <-u.done:
		return ErrRepositoryIsClosing
	case u.deletionQueue <- ud:
		return nil
	}
}

func (u *urlMemoryRepository) fromQueueToBuffer(_ context.Context) {
	for i := 0; i < 10; i++ {
		u.wg.Add(1)
		go func() {
			defer u.wg.Done()
			for {
				select {
				case <-u.done:
					log.Println("Exiting")
					return
				case ud, ok := <-u.deletionQueue:
					if !ok {
						return
					}
					//u.mx.Lock()
					for ind, dto := range u.urlStorage {
						if ud.ID == dto.ID && ud.UserID == dto.UserID {
							u.urlStorage = append(u.urlStorage[:ind], u.urlStorage[ind+1:]...)
							dto.Deleted = true
							u.urlStorage = append(u.urlStorage, dto)
						}
					}
					//u.mx.Unlock()
				}
			}
		}()
	}
}

func (u *urlMemoryRepository) RemoveAll(ctx context.Context, removingList []entity.URLDTO) error {
	fmt.Printf("List to delete %v", removingList)

	for _, dto := range removingList {
		for ind, ud := range u.urlStorage {
			if ud.ID == dto.ID && ud.UserID == dto.UserID {
				u.mx.Lock()
				u.urlStorage = append(u.urlStorage[:ind], u.urlStorage[ind+1:]...)
				ud.Deleted = true
				u.urlStorage = append(u.urlStorage, ud)
				u.mx.Unlock()
			}
		}
	}
	//u.fromQueueToBuffer(ctx)
	//for _, ud := range removingList {
	//	err := u.addURLToDeletionQueue(ud)
	//	if err != nil {
	//		return err
	//	}
	//}
	log.Printf("Repo after delete %v", u.urlStorage)
	return nil
}

func (u *urlMemoryRepository) GetAllByUserID(_ context.Context, userID string) ([]entity.FullURL, error) {
	u.mx.Lock()
	defer u.mx.Unlock()

	listFullURL := make([]entity.FullURL, 0)
	for _, ud := range u.urlStorage {
		if ud.UserID == userID {
			fullURL := entity.FullURL{
				OriginalURL: ud.OriginalURL,
				ShortURL:    fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), ud.ID),
			}
			listFullURL = append(listFullURL, fullURL)
		}
	}
	if len(listFullURL) == 0 {
		return nil, fmt.Errorf("user with id %s has not URL's", userID)
	}
	return listFullURL, nil
}

func (u *urlMemoryRepository) Save(_ context.Context, userID, urlID, url string) error {
	u.mx.Lock()
	defer u.mx.Unlock()

	ud := entity.URLDTO{
		ID:          urlID,
		UserID:      userID,
		OriginalURL: url,
		Deleted:     false,
	}
	if exists(u.urlStorage, urlID) {
		return fmt.Errorf("urlservice %s already exist", url)
	}
	u.urlStorage = append(u.urlStorage, ud)

	return nil
}

func (u *urlMemoryRepository) FindByID(_ context.Context, urlID string) (string, error) {
	u.mx.Lock()
	defer u.mx.Unlock()

	var originalURL string

	for _, ud := range u.urlStorage {
		if ud.ID == urlID {
			if ud.Deleted {
				deletedErr := myerrors.NewDeletedError(ud, nil)
				return "", deletedErr
			}
			originalURL = ud.OriginalURL
		}
	}
	if len(originalURL) == 0 {
		return "", fmt.Errorf("urlservice with id %s not found", urlID)
	}
	return originalURL, nil
}

func (u *urlMemoryRepository) Ping() error {
	return nil
}

func New() repositories.URLRepository {
	return &urlMemoryRepository{
		urlStorage:    make([]entity.URLDTO, 0),
		deletionQueue: make(chan entity.URLDTO),
		done:          make(chan struct{}),
	}
}

func exists(urlStorage []entity.URLDTO, urlID string) bool {
	for _, ud := range urlStorage {
		if ud.ID == urlID {
			return true
		}
	}
	return false
}

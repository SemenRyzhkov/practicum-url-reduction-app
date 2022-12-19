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
	mx sync.Mutex
	//commonURLStorage map[string]string
	//urlStorage       map[string]map[string]string
	urlStorage    []entity.URLDTO
	deletionQueue chan entity.URLDTO
	done          chan struct{}
	wg            sync.WaitGroup
	once          sync.Once
}

func (u *urlMemoryRepository) StopWorkerPool() {
	u.once.Do(func() {
		close(u.done)
	})

	close(u.deletionQueue)
	u.wg.Wait()
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
	for i := 0; i < 10; i++ { // создаем 10 горутин-воркеров
		u.wg.Add(1)
		go func() {
			defer u.wg.Done()
			for {
				select {
				case <-u.done:
					log.Println("Exiting")
					return // если поступает, сигнал из канала done, завершаем
				case ud, ok := <-u.deletionQueue: // вычитываем из очереди
					if !ok {
						return
					}
					for _, dto := range u.urlStorage {
						if ud.ID == dto.ID && ud.UserID == dto.UserID {
							dto.Deleted = true
						}
					}
				}
			}
		}()
	}
}

func (u *urlMemoryRepository) RemoveAll(ctx context.Context, removingList []entity.URLDTO) error {
	u.fromQueueToBuffer(ctx)
	for _, ud := range removingList {
		err := u.addURLToDeletionQueue(ud)
		if err != nil {
			return err
		}
	}
	log.Printf("Repo after delete %v", u.urlStorage)
	//d.Stop()
	return nil
}

func (u *urlMemoryRepository) GetAllByUserID(_ context.Context, userID string) ([]entity.FullURL, error) {
	u.mx.Lock()
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
	u.mx.Unlock()
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
		Deleted:     true,
	}
	if exists(u.urlStorage, urlID) {
		return fmt.Errorf("urlservice %s already exist", url)
	}
	u.urlStorage = append(u.urlStorage, ud)
	return nil
}

func (u *urlMemoryRepository) FindByID(_ context.Context, urlID string) (string, error) {
	u.mx.Lock()
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
	u.mx.Unlock()
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
		//commonURLStorage: make(map[string]string),
		//urlStorage:       make(map[string]map[string]string),
		urlStorage:    make([]entity.URLDTO, 0),
		deletionQueue: make(chan entity.URLDTO),
		done:          make(chan struct{}),
	}
}

//func exists(urlStorage map[string]string, urlID string) bool {
//	_, ok := urlStorage[urlID]
//	return ok
//}

func exists(urlStorage []entity.URLDTO, urlID string) bool {
	for _, ud := range urlStorage {
		if ud.ID == urlID {
			return true
		}
	}
	return false
}

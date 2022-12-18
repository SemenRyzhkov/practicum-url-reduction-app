package indatabase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/omeid/pgerror"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity/myerrors"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
)

var (
	_                      repositories.URLRepository = &dbURLRepository{}
	ErrRepositoryIsClosing                            = errors.New("repository is closing")
)

const (
	initDBQuery = "" +
		"CREATE TABLE IF NOT EXISTS public.urls (" +
		"id varchar(45) primary key, " +
		"original_url text, " +
		"user_id varchar(45), " +
		"deleted boolean" +
		")"
	createUserIDIndex = "" +
		"CREATE INDEX IF NOT EXISTS user_id_index " +
		"ON public.urls (user_id)"
	getAllQuery = "" +
		"SELECT id, original_url " +
		"FROM public.urls " +
		"WHERE user_id=$1"
	getURLQuery = "" +
		"SELECT original_url, deleted FROM public.urls " +
		"WHERE id=$1"
	insertURLQuery = "" +
		"INSERT INTO public.urls (id, original_url, user_id, deleted) " +
		"VALUES ($1, $2, $3, $4)"
	deleteQuery = "" +
		"UPDATE public.urls " +
		"SET deleted = $1 " +
		"WHERE id = $2 AND user_id = $3"
)

type dbURLRepository struct {
	db            *sql.DB
	deletionQueue chan entity.URLDTO
	done          chan struct{}
	wg            sync.WaitGroup
	buffer        []entity.URLDTO
	mx            sync.Mutex
	once          sync.Once
}

func (d *dbURLRepository) RemoveAll(_ context.Context, removingList []entity.URLDTO) error {
	d.fromQueueToBuffer()
	for _, ud := range removingList {
		err := d.addURLToDeletionQueue(ud)
		if err != nil {
			return err
		}
	}
	return d.Stop()
}

func (d *dbURLRepository) addURLToDeletionQueue(ud entity.URLDTO) error {
	select {
	case <-d.done:
		return ErrRepositoryIsClosing
	case d.deletionQueue <- ud:
		log.Println(d.deletionQueue)
		return nil
	}
}

func (d *dbURLRepository) fromQueueToBuffer() {
	for i := 0; i < 10; i++ { // создаем 10 горутин-воркеров
		d.wg.Add(1)
		go func() {
			defer d.wg.Done()
			for {
				select {
				case <-d.done:
					log.Println("Exiting")
					return // если поступает, сигнал из канала done, завершаем
				case ud, ok := <-d.deletionQueue: // вычитываем из очереди
					if !ok {
						return
					}
					err := d.AddURLToBuffer(&ud)
					log.Printf("In queque %v", ud)
					if err != nil {
						log.Println(err)
						return
					} // добавляем в буфер
				}
			}
		}()
	}
}

func (d *dbURLRepository) Stop() error {
	d.once.Do(func() {
		close(d.done)
	}) // todo: do it in sync.Once
	close(d.deletionQueue)
	d.wg.Wait()
	d.mx.Lock()
	err := d.Flush()
	d.mx.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (d *dbURLRepository) AddURLToBuffer(u *entity.URLDTO) error {
	log.Printf("Add url to buffer %s", u.ID)
	d.mx.Lock()
	d.buffer = append(d.buffer, *u)
	d.mx.Unlock()
	if cap(d.buffer) == len(d.buffer) {
		d.mx.Lock()
		err := d.Flush()
		d.mx.Unlock()
		if err != nil {
			return errors.New("cannot add records to the database")
		}
	}
	return nil
}

func (d *dbURLRepository) Flush() error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(deleteQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()
	log.Printf("Buffer contains %d elements", len(d.buffer))
	for _, u := range d.buffer {
		if _, err = stmt.Exec(u.Deleted, u.ID, u.UserID); err != nil {
			if err = tx.Rollback(); err != nil {
				log.Fatalf("update drivers: unable to rollback: %v", err)
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("update drivers: unable to commit: %v", err)
		return err
	}

	d.buffer = d.buffer[:0]
	return nil
}

func New(dbAddress string) (repositories.URLRepository, error) {
	db, err := initDB(dbAddress)
	if err != nil {
		return nil, err
	}
	return &dbURLRepository{
		db:            db,
		buffer:        make([]entity.URLDTO, 100),
		deletionQueue: make(chan entity.URLDTO),
		done:          make(chan struct{}),
	}, nil
}

func (d *dbURLRepository) Save(ctx context.Context, userID, urlID, url string) error {
	_, err := d.db.ExecContext(ctx, insertURLQuery, urlID, url, userID, false)
	if err != nil {
		if e := pgerror.UniqueViolation(err); e != nil {
			return myerrors.NewViolationError(fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), urlID), err)
		}
	}

	return nil
}

func (d *dbURLRepository) FindByID(ctx context.Context, urlID string) (string, error) {
	var ud entity.URLDTO
	row := d.db.QueryRowContext(ctx, getURLQuery, urlID)
	err := row.Scan(&ud.OriginalURL, &ud.Deleted)
	log.Printf("ID %s deleted %v", urlID, ud.Deleted)
	if err != nil {
		return "", err
	}

	if ud.Deleted {
		deletedErr := myerrors.NewDeletedError(ud, nil)
		return "", deletedErr
	}

	return ud.OriginalURL, nil
}

func (d *dbURLRepository) GetAllByUserID(ctx context.Context, userID string) ([]entity.FullURL, error) {
	urls := make([]entity.FullURL, 0)

	rows, err := d.db.QueryContext(ctx, getAllQuery, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u entity.FullURL
		err = rows.Scan(&u.ShortURL, &u.OriginalURL)
		if err != nil {
			return nil, err
		}

		urls = append(urls, u)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func (d *dbURLRepository) Ping() error {
	pingErr := d.db.Ping()
	if pingErr != nil {
		return pingErr
	}
	return nil
}

func initDB(dbAddress string) (*sql.DB, error) {
	db, connectionErr := sql.Open("postgres", dbAddress)
	if connectionErr != nil {
		log.Println(connectionErr)
		return nil, connectionErr
	}
	log.Printf("Connect success %s", dbAddress)
	createTableErr := createTableIfNotExists(db)
	if createTableErr != nil {
		log.Println(createTableErr)
		return nil, createTableErr
	}
	return db, nil
}

func createTableIfNotExists(db *sql.DB) error {
	_, createTableErr := db.Exec(initDBQuery)
	if createTableErr != nil {
		return createTableErr
	}
	_, createIndexErr := db.Exec(createUserIDIndex)
	if createIndexErr != nil {
		return createIndexErr
	}
	return nil
}

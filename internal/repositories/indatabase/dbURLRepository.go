package indatabase

import (
	"context"
	"database/sql"
	"log"

	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/entity"
	"github.com/SemenRyzhkov/practicum-url-reduction-app/internal/repositories"
)

var _ repositories.URLRepository = &dbURLRepository{}

const (
	initDBQuery = "" +
		"CREATE TABLE IF NOT EXISTS public.urls (" +
		"id varchar(45) primary key, " +
		"original_url text, " +
		"user_id varchar(45))"
	getAllQuery = "" +
		"SELECT id, original_url " +
		"FROM public.urls " +
		"WHERE user_id=$1"
	getReduceURLQuery = "" +
		"SELECT original_url FROM public.urls " +
		"WHERE id=$1"
	insertURLQuery = "" +
		"INSERT INTO public.urls (id, original_url, user_id) " +
		"VALUES ($1, $2, $3)"
)

type dbURLRepository struct {
	db *sql.DB
}

func (d dbURLRepository) Ping() error {
	pingErr := d.db.Ping()
	if pingErr != nil {
		return pingErr
	}
	return nil
}

func New(dbAddress string) repositories.URLRepository {
	return &dbURLRepository{
		db: initDB(dbAddress),
	}
}

func (d dbURLRepository) Save(ctx context.Context, userID, urlID, url string) error {
	_, err := d.db.ExecContext(ctx, insertURLQuery, urlID, url, userID)
	if err != nil {
		return err
	}

	return nil
}

func (d dbURLRepository) FindByID(ctx context.Context, urlID string) (string, error) {
	var originalURL string
	row := d.db.QueryRowContext(ctx, getReduceURLQuery, urlID)
	err := row.Scan(&originalURL)
	if err != nil {
		return "", err
	}
	return originalURL, nil
}

func (d dbURLRepository) GetAllByUserID(ctx context.Context, userID string) ([]entity.FullURL, error) {
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

func initDB(dbAddress string) *sql.DB {
	db, connectionErr := sql.Open("postgres", dbAddress)
	if connectionErr != nil {
		log.Fatal(connectionErr)
	}

	createTableIfNotExists(db)
	return db
}

func createTableIfNotExists(db *sql.DB) {
	_, err := db.Exec(initDBQuery)

	if err != nil {
		log.Fatal(err)
	}
}

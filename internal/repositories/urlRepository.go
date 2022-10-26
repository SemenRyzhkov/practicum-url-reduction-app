package repositories

type UrlRepository interface {
	Save(urlId, url string) error
	FindById(urlId string) (string, error)
}

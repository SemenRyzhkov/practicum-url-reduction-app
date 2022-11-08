package repositories

type URLRepository interface {
	Save(urlID, url string) error
	FindByID(urlID string) (string, error)
}

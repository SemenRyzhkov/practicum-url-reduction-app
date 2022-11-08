package service

type URLService interface {
	ReduceAndSaveURL(url string) (string, error)
	GetURLByID(urlID string) (string, error)
}

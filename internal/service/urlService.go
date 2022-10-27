package service

type UrlService interface {
	ReduceAndSaveUrl(url string) (string, error)
	GetUrlById(urlId string) (string, error)
}

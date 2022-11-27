package entity

type FullURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

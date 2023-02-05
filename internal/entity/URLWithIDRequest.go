package entity

// URLWithIDRequest запрос с урлом для сокращения c ID
type URLWithIDRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

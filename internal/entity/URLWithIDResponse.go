package entity

// URLWithIDResponse ответ сервера с сокращенным урлом, cодержащим ID
type URLWithIDResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

package models

type JSONRequest struct {
	URL string `json:"url"`
}

type URLBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

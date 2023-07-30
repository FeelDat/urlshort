package models

type JSONResponse struct {
	Result string `json:"result"`
}

type URLRBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type UsersURLS struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

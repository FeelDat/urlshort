package models

type JsonResponse struct {
	Result string `json:"result"`
}

type URLRBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

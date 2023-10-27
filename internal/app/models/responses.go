// Package models provides data structures and types that are used across the application.
package models

// JSONResponse represents a general response structure containing a result.
//
// It is typically used to convey the result or status of a processed request.
type JSONResponse struct {
	// Result holds the result string of the processed request.
	Result string `json:"result"`
}

// URLRBatchResponse represents a structure for a batch URL response.
//
// It is used to send back shortened URLs with associated correlation IDs after processing.
type URLRBatchResponse struct {
	// CorrelationID is the unique identifier corresponding to each URL in a batch.
	CorrelationID string `json:"correlation_id"`
	// ShortURL is the shortened version of the original URL.
	ShortURL string `json:"short_url"`
}

// UsersURLS represents the structure of URLs associated with a user.
//
// It contains both the original and the shortened version of a URL.
type UsersURLS struct {
	// ShortURL is the shortened version of the original URL.
	ShortURL string `json:"short_url"`
	// OriginalURL is the actual link provided by the user.
	OriginalURL string `json:"original_url"`
}

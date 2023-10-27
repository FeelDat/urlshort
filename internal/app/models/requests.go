// Package models provides data structures and types that are used across the application.
package models

// JSONRequest represents a structure for a single URL request.
//
// It is used to capture incoming request data where a client sends a URL to be processed.
type JSONRequest struct {
	// URL represents the link that needs to be processed.
	URL string `json:"url"`
}

// URLBatchRequest represents a structure for a batch URL request.
//
// It is used when clients send a batch of URLs with associated correlation IDs
// to be processed together.
type URLBatchRequest struct {
	// CorrelationID uniquely identifies each URL in a batch.
	CorrelationID string `json:"correlation_id"`
	// OriginalURL is the actual link that needs to be processed.
	OriginalURL string `json:"original_url"`
}

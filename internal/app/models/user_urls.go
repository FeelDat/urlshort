// Package models provides data structures and types that are used across the application.
package models

// URL represents the structure of a shortened URL within the system.
//
// It contains the original URL, its shortened counterpart, and the associated user's ID.
type URL struct {
	// ShortURL is the shortened version of the original URL.
	ShortURL string

	// OriginalURL is the actual link provided by the user.
	OriginalURL string

	// UserID is the identifier for the user who shortened the URL.
	UserID string
}

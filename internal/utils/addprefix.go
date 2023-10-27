// Package utils contains utility functions used throughout the application.
package utils

import (
	"net/url"
)

// AddPrefix ensures that the provided address (addr) contains a scheme (e.g., "http").
// If the address lacks a scheme, "http" is added by default.
//
// Parameters:
//   - addr: The URL or address to check and potentially prefix.
//
// Returns:
//   - A string representing the address with a scheme.
//   - An error if parsing the address fails.
func AddPrefix(addr string) (string, error) {

	v, err := url.Parse(addr)
	if err != nil {
		return "", err
	}

	if v.Scheme == "" {
		v.Scheme = "http"
	}

	return v.String(), nil
}

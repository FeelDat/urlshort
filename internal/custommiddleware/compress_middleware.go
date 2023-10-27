// Package custommiddleware provides custom middlewares for use in an HTTP server.
package custommiddleware

import (
	"compress/gzip"
	"net/http"
)

// CompressMiddleware is a struct responsible for decompressing request
// bodies that are compressed using the gzip compression algorithm.
type CompressMiddleware struct{}

// NewCompressMiddleware initializes and returns an instance of CompressMiddleware.
func NewCompressMiddleware() *CompressMiddleware {
	return &CompressMiddleware{}
}

// CompressMiddleware is a middleware function that checks if the request's
// Content-Encoding header is set to "gzip". If it is, the middleware decompresses
// the request body using the gzip algorithm before passing the request
// to the next handler in the chain.
func (m *CompressMiddleware) CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = gz
			defer gz.Close()
		}
		next.ServeHTTP(w, r)
	})
}

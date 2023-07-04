package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type CompressMiddleware struct{}

func NewCompressMiddleware() *CompressMiddleware {
	return &CompressMiddleware{}
}

func (m *CompressMiddleware) CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "Invalid gzip content", http.StatusBadRequest)
			return
		}
		r.Body = reader
		next.ServeHTTP(w, r)

		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" && contentType != "text/html" {
			return
		}
		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Close()

		grw := &gzipResponseWriter{
			ResponseWriter: w,
			Writer:         gzipWriter,
		}

		next.ServeHTTP(grw, r)

	})
}

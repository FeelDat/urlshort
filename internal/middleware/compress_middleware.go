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

		acceptEncoding := r.Header.Get("Accept-Encoding")
		if !strings.Contains(acceptEncoding, "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" && contentType != "text/html" {
			next.ServeHTTP(w, r)
			return
		}

		// Устанавливаем заголовок Content-Encoding в gzip
		w.Header().Set("Content-Encoding", "gzip")

		// Создаем gzip.Writer для сжатия ответа
		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Close()

		// Создаем ResponseWriter, который декомпрессирует и записывает данные в gzip.Writer
		grw := &gzipResponseWriter{
			ResponseWriter: w,
			Writer:         gzipWriter,
		}

		// Передаем обработку запроса следующему обработчику с использованием gzipResponseWriter
		next.ServeHTTP(grw, r)
	})
}

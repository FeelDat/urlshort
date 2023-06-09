package custommiddleware

import (
	"compress/gzip"
	"net/http"
)

type CompressMiddleware struct{}

func NewCompressMiddleware() *CompressMiddleware {
	return &CompressMiddleware{}
}

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

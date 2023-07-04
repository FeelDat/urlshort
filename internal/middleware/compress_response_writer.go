package middleware

import (
	"io"
	"net/http"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (grw *gzipResponseWriter) Write(b []byte) (int, error) {
	contentType := grw.Header().Get("Content-Type")
	if contentType != "application/json" && contentType != "text/html" {
		grw.Write(b)
	}
	return grw.Writer.Write(b)
}

package middleware

import (
	"io"
	"net/http"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func (grw *gzipResponseWriter) Write(b []byte) (int, error) {
	return grw.Writer.Write(b)
}

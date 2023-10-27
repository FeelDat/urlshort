// Package custommiddleware provides utility structures to aid in
// logging HTTP response data, such as the size and status code.
package custommiddleware

import "net/http"

// responseData holds data about the HTTP response. This includes
// the status code and size of the response.
type responseData struct {
	// status represents the HTTP status code of the response.
	status int
	// size represents the size of the response in bytes.
	size int
}

// loggingResponseWriter wraps an http.ResponseWriter to capture
// information about the HTTP response. It is used to intercept
// response writes and captures the status code and size.
type loggingResponseWriter struct {
	http.ResponseWriter
	// responseData holds the captured status and size information.
	responseData *responseData
}

// Write writes the data to the underlying http.ResponseWriter and
// captures the size of the data written.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // capture the size
	return size, err
}

// WriteHeader captures the status code and then writes it to the
// underlying http.ResponseWriter.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // capture the status code
}

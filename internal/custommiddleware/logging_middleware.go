// Package custommiddleware provides middleware functionality
// for logging HTTP requests.
package custommiddleware

import (
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

// LoggerMiddleware wraps a zap.SugaredLogger to provide logging functionality.
// It logs information about each HTTP request, including the URI, method, status,
// duration, and size of the response.
type LoggerMiddleware struct {
	// logger is an instance of a zap.SugaredLogger that is used to
	// log information about HTTP requests.
	logger *zap.SugaredLogger
	pool   sync.Pool
}

// NewLoggerMiddleware initializes and returns a new LoggerMiddleware instance.
// It requires a zap.SugaredLogger to be passed in, which it uses for logging.
func NewLoggerMiddleware(logger *zap.SugaredLogger) *LoggerMiddleware {
	return &LoggerMiddleware{
		logger: logger,
		pool: sync.Pool{
			New: func() interface{} {
				return &loggingResponseWriter{
					responseData: &responseData{},
				}
			},
		},
	}
}

// LoggerMiddleware is a middleware handler that logs information about each HTTP request.
// It captures the start time of the request, calculates the duration of the request,
// and logs the URI, method, status, duration, and size using the logger from the LoggerMiddleware struct.
func (m *LoggerMiddleware) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := m.pool.Get().(*loggingResponseWriter)
		lw.ResponseWriter = w
		lw.responseData.status = 0
		lw.responseData.size = 0

		next.ServeHTTP(lw, r)

		duration := time.Since(start)

		m.logger.Info(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", lw.responseData.status,
			"duration", duration,
			"size", lw.responseData.size,
		)

		m.pool.Put(lw)
	})
}

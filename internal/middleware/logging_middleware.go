package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Middleware struct {
	logger *zap.SugaredLogger
}

func NewLoggerMiddleware(logger *zap.SugaredLogger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

func (m *Middleware) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		m.logger.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	})
}

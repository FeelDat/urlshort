package custommiddleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type LoggerMiddleware struct {
	logger *zap.SugaredLogger
}

func NewLoggerMiddleware(logger *zap.SugaredLogger) *LoggerMiddleware {
	return &LoggerMiddleware{
		logger: logger,
	}
}

func (m *LoggerMiddleware) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rd := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   rd,
		}

		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		m.logger.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", rd.status,
			"duration", duration,
			"size", rd.size,
		)
	})
}

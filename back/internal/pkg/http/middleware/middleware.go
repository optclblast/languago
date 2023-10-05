package middleware

import (
	"languago/internal/pkg/logger"
	"net/http"
	"time"
)

type middleware struct {
	log logger.Logger
}

func NewMiddleware(log logger.Logger) *middleware {
	return &middleware{
		log: log,
	}
}

func (m *middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO implement auth middleware
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) RequestValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO implement request validation middleware
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.log.Info("logging middleware",
			logger.LogFields{
				"datetime":     time.Now(),
				"scheme":       r.URL.Scheme,
				"method":       r.Method,
				"path":         r.URL.Path,
				"remote_addr":  r.RemoteAddr,
				"host":         r.Host,
				"user_agent":   r.UserAgent(),
				"referer":      r.Referer(),
				"request_id":   r.Header.Get("X-Request-ID"),
				"content_type": r.Header.Get("Content-Type"),
			},
		)
		next.ServeHTTP(w, r)
	})
}

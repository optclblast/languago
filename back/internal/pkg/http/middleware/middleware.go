package middleware

import (
	"encoding/json"
	"languago/internal/pkg/http/headers"
	"languago/internal/pkg/logger"
	"net/http"
	"strings"
	"time"
)

type middleware struct {
	log    logger.Logger
	closed bool
}

func NewMiddleware(log logger.Logger) *middleware {
	return &middleware{log: log}
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

func (m *middleware) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				m.log.Warn("recovered from panic", logger.LogFields{
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

				jsonBody, _ := json.Marshal(map[string]string{
					"error": "Internal server error",
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func DoAuth(r *http.Request) bool {
	if (r.Header.Get(headers.H_SIGN_IN) != "" || r.Header.Get(headers.H_SIGN_UP) != "") &&
		(r.Method == http.MethodPost && (strings.Contains(r.RequestURI, "signin") ||
			strings.Contains(r.RequestURI, "signup"))) {
		return false

	}

	return true
}

package middleware

import (
	"encoding/json"
	"languago/internal/pkg/ctxtools"
	"languago/internal/pkg/http/headers"
	"languago/internal/pkg/logger"
	"net/http"
	"strings"
	"time"
)

type middleware struct {
	log logger.Logger
}

func NewMiddleware(log logger.Logger) *middleware {
	return &middleware{log: log}
}

func (m *middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !doAuth(r) {
			m.log.Warn(
				"authorization request (sign up)",
				logger.LogFields{
					"datetime":    time.Now(),
					"request_id":  ctxtools.RequestId(r.Context()),
					"remote_addr": r.RemoteAddr,
					"host":        r.Host,
					"user_agent":  r.UserAgent(),
					"referer":     r.Referer(),
				},
			)
		}

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
		m.logRequest(r, "log", false)
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				m.logRequest(r, "recover", true)

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

func doAuth(r *http.Request) bool {
	if r.Header.Get(headers.H_SIGN_UP) != "" &&
		(r.Method == http.MethodPost &&
			strings.Contains(r.RequestURI, "signup")) {
		return false
	}

	return true
}

func (m *middleware) logRequest(r *http.Request, mw string, err bool) {
	fields := logger.LogFields{
		"datetime":     time.Now(),
		"request_id":   ctxtools.RequestId(r.Context()),
		"scheme":       r.URL.Scheme,
		"method":       r.Method,
		"path":         r.URL.Path,
		"remote_addr":  r.RemoteAddr,
		"host":         r.Host,
		"user_agent":   r.UserAgent(),
		"referer":      r.Referer(),
		"content_type": r.Header.Get("Content-Type"),
	}

	if err {
		m.log.Error(mw, fields)
	} else {
		m.log.Info(mw, fields)
	}
}

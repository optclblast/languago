package middleware

import (
	"context"
	"encoding/json"
	"languago/internal/pkg/auth"
	"languago/internal/pkg/ctxtools"
	"languago/internal/pkg/logger"
	"net/http"
	"time"

	errors2 "languago/internal/pkg/errors"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const (
	H_SIGN_IN       = "L-Sign-In"
	H_SIGN_UP       = "L-Sign-Up"
	H_Authorization = "Authorization"
)

type middleware struct {
	log  logger.Logger
	auth auth.Authorizer
}

func NewMiddleware(log logger.Logger, auth auth.Authorizer) *middleware {
	return &middleware{
		log:  log,
		auth: auth,
	}
}

func (m *middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !doAuth(r) {
			m.log.Warn(
				"[ SIGN_UP_REQUEST ]",
				logger.LogFields{
					"datetime":    time.Now(),
					"request_id":  ctxtools.RequestId(r.Context()),
					"remote_addr": r.RemoteAddr,
					"host":        r.Host,
					"user_agent":  r.UserAgent(),
					"referer":     r.Referer(),
				},
			)

			userID := uuid.New()

			ctxR := r.WithContext(context.WithValue(r.Context(), ctxtools.UserIDCtxKey, userID))

			token, err := m.auth.CreateToken(auth.ClaimJWTParams{
				UserId: userID.String(),
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ctxR.Header.Add(H_Authorization, token)

			next.ServeHTTP(w, ctxR)
		} else {
			tokenStr := r.Header.Get(H_Authorization)
			if tokenStr == "" {
				m.log.Error("error auth", logger.LogFields{
					"datetime":    time.Now(),
					"request_id":  ctxtools.RequestId(r.Context()),
					"remote_addr": r.RemoteAddr,
					"host":        r.Host,
					"user_agent":  r.UserAgent(),
					"referer":     r.Referer(),
				})
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			token, err := jwt.ParseWithClaims(tokenStr, make(jwt.MapClaims), func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
					return nil, errors2.ErrInvalidToken
				}
				return m.auth.Secret(), nil
			})
			if err != nil {
				m.log.Error("error parse token", logger.LogFields{
					"datetime":    time.Now(),
					"request_id":  ctxtools.RequestId(r.Context()),
					"remote_addr": r.RemoteAddr,
					"host":        r.Host,
					"user_agent":  r.UserAgent(),
					"referer":     r.Referer(),
					"error":       err,
				})
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			user, err := m.auth.Authorize(token)
			if err != nil {
				m.log.Error("error auth", logger.LogFields{
					"datetime":    time.Now(),
					"request_id":  ctxtools.RequestId(r.Context()),
					"remote_addr": r.RemoteAddr,
					"host":        r.Host,
					"user_agent":  r.UserAgent(),
					"referer":     r.Referer(),
					"error":       err,
				})
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ctxR := r.WithContext(context.WithValue(r.Context(), ctxtools.UserCtxKey, user))

			next.ServeHTTP(w, ctxR)
		}
	})
}

func (m *middleware) RequestValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO implement request validation middleware
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) Options(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
	})
}

func (m *middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logRequest(r, "[ REQUEST_LOG ]", false)
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				m.log.Error("error auth", logger.LogFields{
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
					"error":        err,
				})

				jsonBody, _ := json.Marshal(map[string]string{
					"error": "Internal server error",
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write(jsonBody)
				if err != nil {
					m.log.Error("error write to connection", logger.LogFields{
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
						"error":        err,
					})
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}

var testKey string = "1234"

func doAuth(r *http.Request) bool {
	if r.Header[H_SIGN_UP][0] == testKey && r.Method == http.MethodPost {
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
		"body":         r.Body,
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

package middleware

import (
	"context"
	"encoding/json"
	"languago/infrastructure/logger"
	"languago/pkg/auth"
	"languago/pkg/ctxtools"
	"net/http"
	"time"

	errors2 "languago/pkg/errors"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const (
	H_SIGN_IN       = "L-Sign-In"
	H_SIGN_UP       = "L-Sign-Up"
	H_Authorization = "Authorization"
)

type middleware struct {
	log zerolog.Logger
}

func NewMiddleware(log zerolog.Logger) *middleware {
	return &middleware{log: log}
}

func (m *middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if !doAuth(r) {
			m.log.Warn().Msgf(
				`[ SIGN_UP_REQUEST ] 
				datetime: %v 
				request_id: %v 
				remote_addr: %v 
				host: %v 
				user_agent: %v 
				referer: %v`,
				time.Now(),
				ctxtools.RequestId(r.Context()),
				r.RemoteAddr,
				r.Host,
				r.UserAgent(),
				r.Referer(),
			)

			userID := uuid.New()

			ctxR := r.WithContext(
				context.WithValue(r.Context(), ctxtools.UserIDCtxKey, userID),
			)

			token, err := auth.CreateToken(auth.ClaimJWTParams{
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
				m.log.Error().Msgf("error auth: %v", logger.LogFields{
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

			token, err := jwt.ParseWithClaims(tokenStr, make(jwt.MapClaims),
				func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
						return nil, errors2.ErrInvalidToken
					}
					return auth.Secret(), nil
				})
			if err != nil {
				m.log.Error().Msgf("error parse token: %v", logger.LogFields{
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

			user, err := auth.Authorize(token)
			if err != nil {
				m.log.Error().Msgf("error auth: %v", logger.LogFields{
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

func (m *middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logRequest(r, "[ REQUEST_LOG ]")
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				// m.log.Error().Msgf(
				// 	"fatal error: %v",
				// 	"datetime", time.Now(),
				// 	"request_id", ctxtools.RequestId(r.Context()),
				// 	"scheme", r.URL.Scheme,
				// 	"method", r.Method,
				// 	"path", r.URL.Path,
				// 	"remote_addr", r.RemoteAddr,
				// 	"host", r.Host,
				// 	"user_agent", r.UserAgent(),
				// 	"referer", r.Referer(),
				// 	"content_type", r.Header.Get("Content-Type"),
				// 	"error", err,
				// )

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

var testKey string = "1234"

func doAuth(r *http.Request) bool {
	if r.Method == http.MethodOptions {
		return false
	}

	// if r.Header[H_SIGN_UP] == testKey && r.Method == http.MethodPost {
	// 	return false
	// }

	//return true
	return false
}

func (m *middleware) logRequest(r *http.Request, mw string) {
	m.log.Info().Msgf(
		"%v", []interface{}{
			mw,
			"datetime", time.Now(),
			"request_id", ctxtools.RequestId(r.Context()),
			"scheme", r.URL.Scheme,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"host", r.Host,
			"user_agent", r.UserAgent(),
			"referer", r.Referer(),
			"content_type", r.Header.Get("Content-Type"),
		},
	)
}

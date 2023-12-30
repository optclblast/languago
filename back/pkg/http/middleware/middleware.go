package middleware

import (
	"context"
	"encoding/json"
	"languago/pkg/auth"
	"languago/pkg/ctxtools"
	"net/http"

	errors2 "languago/pkg/errors"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	H_SIGN_IN       = "L-Sign-In"
	H_SIGN_UP       = "L-Sign-Up"
	H_Authorization = "Authorization"
)

type middleware struct {
	log *logrus.Logger
}

func NewMiddleware(log *logrus.Logger) *middleware {
	return &middleware{log: log}
}

func (m *middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if !doAuth(r) {
			reqID, _ := ctxtools.RequestId(r.Context())

			userID := uuid.New()

			ctxR := r.WithContext(
				context.WithValue(r.Context(), ctxtools.UserIDCtxKey, userID),
			)

			token, err := auth.CreateToken(auth.ClaimJWTParams{
				UserId: userID.String(),
			})
			if err != nil {
				m.log.Errorf("reqID: %s, error: %s", reqID, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ctxR = ctxR.WithContext(
				context.WithValue(ctxR.Context(), ctxtools.TokenCtxKey, token),
			)

			ctxR.Header.Add(H_Authorization, token)

			next.ServeHTTP(w, ctxR)
		} else {
			tokenStr := r.Header.Get(H_Authorization)
			reqID, _ := ctxtools.RequestId(r.Context())
			if tokenStr == "" {
				m.log.Errorf("reqID: %s, error: empty token", reqID)
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
				reqID, _ := ctxtools.RequestId(r.Context())
				m.log.Errorf("reqID: %s, error: %s", reqID, err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			user, err := auth.Authorize(token)
			if err != nil {
				reqID, _ := ctxtools.RequestId(r.Context())
				m.log.Errorf("reqID: %s, error: %s", reqID, err.Error())
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

func (m *middleware) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				reqID, _ := ctxtools.RequestId(r.Context())
				m.log.Errorf("reqID: %s, error: %v", reqID, err)

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

	if len(r.Header[H_SIGN_UP]) > 0 {
		if r.Header[H_SIGN_UP][0] == testKey && r.Method == http.MethodPost {
			return false
		}
	}

	//return true
	return false
}

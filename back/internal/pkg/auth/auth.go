package auth

import (
	"context"
	"fmt"
	"languago/internal/pkg/logger"
	"languago/internal/pkg/models"
	"languago/internal/pkg/repository"
	"time"

	errors2 "languago/internal/pkg/errors"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Authorizer interface {
	Authorize(token *jwt.Token) (*models.User, error)
	CreateToken(c ClaimJWTParams) (string, error)
	Secret() []byte
}

type authorizer struct {
	log         logger.Logger
	secret      []byte
	userStorage repository.UserRepository
}

func NewAuthorizer(log logger.Logger, userStorage repository.UserRepository, secret []byte) Authorizer {
	return &authorizer{
		log:         log,
		secret:      secret,
		userStorage: userStorage,
	}
}

func (a *authorizer) Authorize(token *jwt.Token) (*models.User, error) {
	if err := token.Claims.Valid(); err != nil {
		a.log.Warn("invalid token claims", logger.LogFields{
			"module": "authorizer",
			"error":  err,
		})
		return nil, errors2.ErrInvalidToken
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		a.log.Warn("invalid token claims", nil)
		return nil, errors2.ErrInvalidToken
	}

	userIDstr, ok := payload["sub"].(string)
	if !ok || userIDstr == "" {
		a.log.Warn("invalid token claims", logger.LogFields{
			"module": "authorizer",
			"error":  "invalid sub claim",
			"is_ok":  ok,
		})
		return nil, errors2.ErrInvalidToken
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		a.log.Warn("error parse user_id", logger.LogFields{
			"module": "authorizer",
		})
		return nil, errors2.ErrInvalidToken
	}

	user, err := a.userStorage.SelectUser(ctx, repository.SelectUserParams{
		ID: userID,
	})
	if err != nil {
		a.log.Warn("error select user", logger.LogFields{
			"error": err,
		})
		return nil, errors2.ErrUnauthorized
	}

	a.log.Info("[ AUTHORIZE ] user authorized", logger.LogFields{
		"module":  "authorizer",
		"user_id": userID,
	})
	return user.ToModel(), nil
}

type ClaimJWTParams struct {
	UserId string
}

func (a *authorizer) CreateToken(c ClaimJWTParams) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Subject:   c.UserId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(a.secret)
	if err != nil {
		a.log.Error("error sign token", logger.LogFields{
			"module": "authorizer",
			"error":  err,
		})
		return "", fmt.Errorf("error sign token: %w", err)
	}

	return signed, nil
}

func (a *authorizer) Secret() []byte {
	return a.secret
}

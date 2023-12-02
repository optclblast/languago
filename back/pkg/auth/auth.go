package auth

import (
	"context"
	"fmt"
	"languago/infrastructure/repository"
	"languago/pkg/models"
	"time"

	errors2 "languago/pkg/errors"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Authorizer interface {
	Authorize(token *jwt.Token) (*models.User, error)
	CreateToken(c ClaimJWTParams) (string, error)
	Secret() []byte
}

type authorizer struct {
	log         zerolog.Logger
	secret      []byte
	userStorage repository.UserRepository
}

func NewAuthorizer(log zerolog.Logger, userStorage repository.UserRepository, secret []byte) Authorizer {
	return &authorizer{
		log:         log,
		secret:      secret,
		userStorage: userStorage,
	}
}

func (a *authorizer) Authorize(token *jwt.Token) (*models.User, error) {
	if err := token.Claims.Valid(); err != nil {
		a.log.Warn().Msg(fmt.Sprintf("invalid token claims: %s", err.Error()))
		return nil, errors2.ErrInvalidToken
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		a.log.Warn().Msg("invalid token claims")
		return nil, errors2.ErrInvalidToken
	}

	userIDstr, ok := payload["sub"].(string)
	if !ok || userIDstr == "" {
		a.log.Warn().Msg("invalid sub claim")
		return nil, errors2.ErrInvalidToken
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		a.log.Warn().Msg("auth: error parse user_id")
		return nil, errors2.ErrInvalidToken
	}

	user, err := a.userStorage.SelectUser(ctx, repository.SelectUserParams{
		ID: userID,
	})
	if err != nil {
		a.log.Warn().Msg("error select user: " + err.Error())
		return nil, errors2.ErrUnauthorized
	}

	a.log.Info().Msg(fmt.Sprintf("[ AUTHORIZE ] user authorized. user: %s time: %v"+user.Id.String(), time.Now()))
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
		a.log.Error().Msg("error sign token")
		return "", fmt.Errorf("error sign token: %w", err)
	}

	return signed, nil
}

func (a *authorizer) Secret() []byte {
	return a.secret
}

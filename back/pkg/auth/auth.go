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
	"github.com/sirupsen/logrus"
)

type authorizer struct {
	log         logrus.Logger
	secret      []byte
	userStorage repository.UserRepository
}

var _authorizer authorizer

func NewAuthorizer(log *logrus.Logger, userStorage repository.UserRepository, secret []byte) {
	_authorizer = authorizer{
		log:         *log,
		secret:      secret,
		userStorage: userStorage,
	}
}

func Authorize(token *jwt.Token) (*models.User, error) {
	if err := token.Claims.Valid(); err != nil {
		_authorizer.log.Warn(fmt.Sprintf("invalid token claims: %s", err.Error()))
		return nil, errors2.ErrInvalidToken
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		_authorizer.log.Warn("invalid token claims")
		return nil, errors2.ErrInvalidToken
	}

	userIDstr, ok := payload["sub"].(string)
	if !ok || userIDstr == "" {
		_authorizer.log.Warn("invalid sub claim")
		return nil, errors2.ErrInvalidToken
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		_authorizer.log.Warn("auth: error parse user_id")
		return nil, errors2.ErrInvalidToken
	}

	user, err := _authorizer.userStorage.SelectUser(ctx, repository.SelectUserParams{
		ID: userID,
	})
	if err != nil {
		_authorizer.log.Warn("error select user: " + err.Error())
		return nil, errors2.ErrUnauthorized
	}

	_authorizer.log.Info(fmt.Sprintf("[ AUTHORIZE ] user authorized. user: %s time: %v"+user.Id.String(), time.Now()))
	return user.ToModel(), nil
}

type ClaimJWTParams struct {
	UserId string
}

func CreateToken(c ClaimJWTParams) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Subject:   c.UserId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(_authorizer.secret)
	if err != nil {
		_authorizer.log.Error("error sign token")
		return "", fmt.Errorf("error sign token: %w", err)
	}

	return signed, nil
}

func Secret() []byte {
	return _authorizer.secret
}

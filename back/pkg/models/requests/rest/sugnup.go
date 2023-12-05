package rest

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	ID    uuid.UUID  `json:"id"`
	Token *jwt.Token `json:"token"`
}

type SignInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token     string `json:"token"`
	ExpitedAt uint64 `json:expired_at`
}

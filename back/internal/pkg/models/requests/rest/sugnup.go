package rest

import (
	"languago/internal/pkg/auth"

	"github.com/google/uuid"
)

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	ID    uuid.UUID  `json:"id"`
	Token auth.Token `json:"token"`
}
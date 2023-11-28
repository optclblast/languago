package rest

import "github.com/google/uuid"

type (
	CreateUserRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	CreateUserResponse struct {
		ID    uuid.UUID `json:"id"`
		Login string    `json:"login"`
	}

	GetUserRequest  struct{}
	GetUserResponse struct{}

	EditUserRequest  struct{}
	EditUserResponse struct{}

	DeleteUserRequest  struct{}
	DeleteUserResponse struct{}
)

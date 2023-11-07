package users

import (
	"context"
	"fmt"
	"languago/internal/pkg/logger"
	"languago/internal/pkg/models/requests/rest"
	"languago/internal/pkg/repository"

	"github.com/google/uuid"
)

type UsersController interface {
	CreateUser(ctx context.Context, req *rest.SignUpRequest) error
	GetUser(ctx context.Context, req *rest.GetUserRequest) (*rest.GetUserResponse, error)
	DeleteUser(ctx context.Context, req *rest.DeleteUserRequest) error
	EditUser(ctx context.Context, req *rest.EditUserRequest) (*rest.EditUserResponse, error)
}

type usersController struct {
	log     logger.Logger
	storage repository.DatabaseInteractor
}

func NewUsersController(
	log logger.Logger,
	storage repository.DatabaseInteractor,
) UsersController {
	return &usersController{
		log:     log,
		storage: storage,
	}
}

func (c *usersController) CreateUser(ctx context.Context, req *rest.SignUpRequest) error {
	userID := uuid.New()
	err := c.storage.Database().CreateUser(ctx, repository.CreateUserParams{
		ID:       userID,
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		return fmt.Errorf("error create new user: %w", err)
	}

	return nil
}

// todo
func (c *usersController) GetUser(ctx context.Context, req *rest.GetUserRequest) (*rest.GetUserResponse, error) {
	return nil, nil
}

// todo
func (c *usersController) DeleteUser(ctx context.Context, req *rest.DeleteUserRequest) error {
	return nil
}

// todo
func (c *usersController) EditUser(ctx context.Context, req *rest.EditUserRequest) (*rest.EditUserResponse, error) {
	return nil, nil
}

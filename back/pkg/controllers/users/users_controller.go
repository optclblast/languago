package users

import (
	"context"
	"fmt"
	"languago/infrastructure/repository"
	"languago/pkg/ctxtools"
	"languago/pkg/models/requests/rest"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type UsersController interface {
	CreateUser(ctx context.Context, req *rest.SignUpRequest) error
	GetUser(ctx context.Context, req *rest.GetUserRequest) (*rest.GetUserResponse, error)
	DeleteUser(ctx context.Context, req *rest.DeleteUserRequest) error
	EditUser(ctx context.Context, req *rest.EditUserRequest) (*rest.EditUserResponse, error)
}

type usersController struct {
	log     zerolog.Logger
	storage repository.DatabaseInteractor
}

func NewUsersController(
	log zerolog.Logger,
	storage repository.DatabaseInteractor,
) UsersController {
	return &usersController{
		log:     log,
		storage: storage,
	}
}

func (c *usersController) CreateUser(ctx context.Context, req *rest.SignUpRequest) error {
	userID := ctxtools.UserID(ctx)
	if userID == uuid.Nil {
		return fmt.Errorf("error fetch user id from context")
	}

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

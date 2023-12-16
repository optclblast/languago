package ctxtools

import (
	"context"
	"database/sql"
	"fmt"
	"languago/pkg/models"

	chi "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type ContextKey string

var (
	UserIDCtxKey         ContextKey = "user_id"
	UserCtxKey           ContextKey = "user"
	TokenCtxKey          ContextKey = "token"
	IsolationLevelCtxKey ContextKey = "isolation_level"
)

func Token(ctx context.Context) (string, error) {
	if v := ctx.Value(TokenCtxKey); v != nil {
		if token, ok := v.(string); ok {
			return token, nil
		}
	}
	return "", fmt.Errorf("token not passed through context")
}

func IsolationLevel(ctx context.Context) (sql.IsolationLevel, error) {
	if v := ctx.Value(chi.RequestIDKey); v != nil {
		if level, ok := v.(sql.IsolationLevel); ok {
			return level, nil
		}
	}
	return -1, fmt.Errorf("isolatin level not passed through context")
}

func RequestId(ctx context.Context) (string, error) {
	if v := ctx.Value(chi.RequestIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s, nil
		}
	}

	return "", fmt.Errorf("request id not passed through context")
}

func UserID(ctx context.Context) (uuid.UUID, error) {
	if v := ctx.Value(UserIDCtxKey); v != nil {
		if s, ok := v.(uuid.UUID); ok {
			return s, nil
		}
	}

	return uuid.Nil, fmt.Errorf("user id not passed through context")
}

func User(ctx context.Context) (*models.User, error) {
	if v := ctx.Value(UserCtxKey); v != nil {
		if s, ok := v.(*models.User); ok {
			return s, nil
		}
	}

	return nil, fmt.Errorf("user not passed through context")
}

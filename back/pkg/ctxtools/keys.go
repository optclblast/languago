package ctxtools

import (
	"context"
	"database/sql"
	"languago/pkg/models"

	chi "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type ContextKey string

var (
	UserIDCtxKey         ContextKey = "user_id"
	UserCtxKey           ContextKey = "user"
	IsolationLevelCtxKey ContextKey = "isolation_level"
)

func IsolationLevel(ctx context.Context) sql.IsolationLevel {
	if v := ctx.Value(chi.RequestIDKey); v != nil {
		if level, ok := v.(sql.IsolationLevel); ok {
			return level
		}
	}
	return -1
}

func RequestId(ctx context.Context) string {
	if v := ctx.Value(chi.RequestIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}

	return ""
}

func UserID(ctx context.Context) uuid.UUID {
	if v := ctx.Value(UserIDCtxKey); v != nil {
		if s, ok := v.(uuid.UUID); ok {
			return s
		}
	}

	return uuid.Nil
}

func User(ctx context.Context) *models.User {
	if v := ctx.Value(UserCtxKey); v != nil {
		if s, ok := v.(*models.User); ok {
			return s
		}
	}

	return nil
}

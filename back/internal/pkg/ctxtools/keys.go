package ctxtools

import (
	"context"

	chi "github.com/go-chi/chi/v5/middleware"
)

type ContextKey string

func RequestId(ctx context.Context) string {
	if v := ctx.Value(chi.RequestIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}

	return ""
}

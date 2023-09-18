package api

import "context"

type (
	Authorizer interface {
		Authorize(ctx context.Context) error
	}

	// DefaultAuthorizer struct {}
	mockAuthorizer struct{}
)

// mock
func (ma *mockAuthorizer) Authorize(ctx context.Context) error {
	return nil
}

// mock

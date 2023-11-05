package errors

import (
	"fmt"
	"languago/internal/pkg/logger"

	"github.com/google/uuid"
)

// TODO errors
const (
	CodeInternalServerError Code = 500
	CodeBadRequest          Code = 400
	CodeNotFound            Code = 404
)

var (
	ErrInternalServerError = New(CodeInternalServerError, "Oops! Something went wrong!")
	ErrNotFound            = New(CodeNotFound, "Not Found")
	ErrValidation          = New(CodeBadRequest, "Validation Error")
	ErrBadRequest          = New(CodeBadRequest, "BadRequest")
)

type Code uint64

type serviceError struct {
	ServiceName string
	ServiceID   uuid.UUID
	Code        Code
	Message     string
	Err         error
}

func (e serviceError) Error() string {
	return fmt.Sprintf(
		"ServiceID: [%v] ServiceName: %s Message: %s Error: %s",
		e.ServiceID,
		e.ServiceName,
		e.Message,
		e.Err.Error(),
	)
}

func New(code Code, msg string, parent ...error) error {
	return nil
}

type ErrorsPersenter interface {
	ServiceError(opts ...ServiceErrorOption) error
	ResponseError(err error) error
}

type errorPresenter struct {
	log logger.Logger
}

func NewErrorPresenter(log logger.Logger) ErrorsPersenter {
	return &errorPresenter{log: log}
}

func (e *errorPresenter) ResponseError(err error) error {
	return nil
}

type ServiceErrorOption func(e *serviceError)

func (e *errorPresenter) ServiceError(opts ...ServiceErrorOption) error {
	var err serviceError
	for _, option := range opts {
		option(&err)
	}

	return err
}

func (e *errorPresenter) mapError(err error) error {
	return nil
}

func ErrorServiceID(serviceID uuid.UUID) ServiceErrorOption {
	return func(e *serviceError) {
		e.ServiceID = serviceID
	}
}

func ErrorServiceName(name string) ServiceErrorOption {
	return func(e *serviceError) {
		e.ServiceName = name
	}
}

func ErrorMessage(msg string) ServiceErrorOption {
	return func(e *serviceError) {
		e.Message = msg
	}
}

func ErrorServiceErr(err error) ServiceErrorOption {
	return func(e *serviceError) {
		e.Err = err
	}
}

package errors

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// TODO errors
const (
	CodeInternalServerError Code = 500
	CodeBadRequest          Code = 400
	CodeNotFound            Code = 404
	CodeUnauthorized        Code = 401
)

var (
	ErrInternalServerError = New(CodeInternalServerError, "Oops! Something went wrong!")
	ErrNotFound            = New(CodeNotFound, "Not Found")
	ErrValidation          = New(CodeBadRequest, "Validation Error")
	ErrBadRequest          = New(CodeBadRequest, "BadRequest")
	ErrInvalidToken        = New(CodeUnauthorized, "Invalid Token")
	ErrUnauthorized        = New(CodeUnauthorized, "Unauthorized")
)

type Code int

type serviceError struct {
	serviceName string
	serviceID   uuid.UUID
	Code        Code   `json:"code"`
	Message     string `json:"message"`
	Err         error  `json:"error,omitempty"`
}

func (e serviceError) Error() string {
	return fmt.Sprintf(
		"ServiceID: [%v] ServiceName: [%s] Message: %s Error: %s",
		e.serviceID,
		e.serviceName,
		e.Message,
		e.Err.Error(),
	)
}

// Returns bare service error. Can be modified using FOP
func New(code Code, msg string, parent ...error) error {
	if parent == nil {
		return serviceError{
			Code:    code,
			Message: msg,
		}
	}

	return serviceError{
		Code:    code,
		Message: msg,
		Err:     errors.Join(parent...),
	}
}

type ErrorsPersenter interface {
	ServiceError(be error, opts ...ServiceErrorOption) error
	ResponseError(err error) error
}

type errorPresenter struct {
	log *logrus.Logger
}

func NewErrorPresenter(log *logrus.Logger) ErrorsPersenter {
	return &errorPresenter{log: log}
}

func (e *errorPresenter) ResponseError(err error) error {
	return e.mapError(err)
}

type ServiceErrorOption func(e *serviceError)

func (e *errorPresenter) ServiceError(be error, opts ...ServiceErrorOption) error {
	var err serviceError
	if serr, ok := be.(serviceError); ok {
		err = serr
	}

	for _, option := range opts {
		option(&err)
	}

	return err
}

func (e *errorPresenter) mapError(err error) error {
	switch {
	case errors.Is(err, ErrBadRequest):
		return ErrBadRequest
	case errors.Is(err, ErrInternalServerError):
		return ErrInternalServerError
	case errors.Is(err, ErrInvalidToken):
		return ErrUnauthorized
	case errors.Is(err, ErrUnauthorized):
		return ErrUnauthorized
	case errors.Is(err, ErrValidation):
		return ErrBadRequest
	default:
		return ErrInternalServerError
	}
}

func ErrorServiceID(serviceID uuid.UUID) ServiceErrorOption {
	return func(e *serviceError) {
		e.serviceID = serviceID
	}
}

func ErrorServiceName(name string) ServiceErrorOption {
	return func(e *serviceError) {
		e.serviceName = name
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

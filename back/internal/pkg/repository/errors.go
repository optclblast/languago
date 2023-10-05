package repository

import (
	"database/sql"

	"github.com/lib/pq"
)

type (
	ErrNotFound struct {
		err error
	}

	ErrInternal struct {
		err error
	}

	ErrAuth struct {
		err error
	}

	ErrChannelAlreadyOpen struct {
		err error
	}

	ErrChannelNotOpen struct {
		err error
	}
)

func (e *ErrNotFound) Error() string {
	return e.err.Error()
}

func (e *ErrInternal) Error() string {
	return e.err.Error()
}

func (e *ErrAuth) Error() string {
	return e.err.Error()
}

func (e *ErrChannelAlreadyOpen) Error() string {
	return e.err.Error()
}

func (e *ErrChannelNotOpen) Error() string {
	return e.err.Error()
}

func properError(err error) error {
	switch {
	case err == nil:
		return nil
	case err == sql.ErrNoRows:
		return &ErrNotFound{err}
	case err == pq.ErrCouldNotDetectUsername ||
		err == pq.ErrSSLKeyUnknownOwnership ||
		err == pq.ErrSSLNotSupported:

		return &ErrAuth{err}
	case err == pq.ErrChannelAlreadyOpen:
		return &ErrChannelAlreadyOpen{err}
	case err == pq.ErrChannelNotOpen:
		return &ErrChannelNotOpen{err}
	default:
		return err
	}
}

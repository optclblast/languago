package repository

import (
	"database/sql"
	errors2 "languago/pkg/errors"

	"github.com/lib/pq"
)

var (
	ErrChannelAlreadyOpen = errors2.New(500, "error channel alreay open", errors2.ErrInternalServerError)
	ErrChannelNotOpen     = errors2.New(500, "error channel not open", errors2.ErrInternalServerError)
	ErrInvalidData        = errors2.New(404, "error invalid data", errors2.ErrValidation)
)

func handleError(err error) error {
	switch {
	case err == nil:
		return nil
	case err == sql.ErrNoRows:
		return errors2.ErrNotFound
	case err == pq.ErrCouldNotDetectUsername ||
		err == pq.ErrSSLKeyUnknownOwnership ||
		err == pq.ErrSSLNotSupported:

		return errors2.ErrBadRequest
	case err == pq.ErrChannelAlreadyOpen:
		return ErrChannelAlreadyOpen
	case err == pq.ErrChannelNotOpen:
		return ErrChannelNotOpen
	default:
		return err
	}
}

package languagoerr

import (
	"fmt"

	"github.com/google/uuid"
)

type (
	FatalErr struct {
		Message     string
		ServiceName string
		ServiceID   uuid.UUID
		Wraps       *FatalErr
	}

	NewFatalErrorParams struct {
		ServiceName string
		ServiceID   uuid.UUID
		Error       error
	}
)

func (e *FatalErr) Error() string {
	var err string = e.Message
	var c *FatalErr = e
	for c.Wraps != nil {
		err = fmt.Sprintf("%s. %s", err, c.Message)
		c = c.Wraps
	}
	return err
}

func NewFatalError(args NewFatalErrorParams, parent ...error) *FatalErr {
	err := FatalErr{
		Message:     args.Error.Error(),
		ServiceName: args.ServiceName,
		ServiceID:   args.ServiceID,
	}

	if len(parent) > 0 {
		var tmpErr *FatalErr = &FatalErr{
			Message: parent[0].Error(),
		}
		for _, e := range parent {
			err := &FatalErr{
				Message: "error: " + e.Error(),
				Wraps:   tmpErr,
			}
			tmpErr = err
		}
		err.Wraps = tmpErr
	}
	return &err
}

package errors

import (
	"github.com/rs/zerolog"
)

type ErrorObservable interface {
	ErrorChannel() chan error
}

type ErrorsObserver interface {
	WatchErrors(ErrorObservable)
}

type errorObserver struct {
	log zerolog.Logger
}

func NewErrorObserver(log zerolog.Logger) ErrorsObserver {
	return &errorObserver{log: log}
}

func (o *errorObserver) WatchErrors(target ErrorObservable) {
	go func(target chan error) {
		err := <-target
		o.log.Error().Msg(err.Error())

	}(target.ErrorChannel())
}

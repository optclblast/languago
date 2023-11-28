package errors

import (
	"languago/infrastructure/logger"
)

type ErrorObservable interface {
	ErrorChannel() chan error
}

type ErrorsObserver interface {
	WatchErrors(ErrorObservable)
}

type errorObserver struct {
	log logger.Logger
}

func NewErrorObserver(log logger.Logger) ErrorsObserver {
	return &errorObserver{log: log}
}

func (o *errorObserver) WatchErrors(target ErrorObservable) {
	go func(target chan error) {
		for {
			select {
			case e := <-target:
				o.log.Error("error: ", logger.LogFieldPair("", e))
			}
		}
	}(target.ErrorChannel())
}

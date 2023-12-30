package errors

import "github.com/sirupsen/logrus"

type ErrorObservable interface {
	ErrorChannel() chan error
}

type ErrorsObserver interface {
	WatchErrors(ErrorObservable)
}

type errorObserver struct {
	log *logrus.Logger
}

func NewErrorObserver(log *logrus.Logger) ErrorsObserver {
	return &errorObserver{log: log}
}

func (o *errorObserver) WatchErrors(target ErrorObservable) {
	go func(target chan error) {
		err := <-target
		o.log.Error(err.Error())

	}(target.ErrorChannel())
}

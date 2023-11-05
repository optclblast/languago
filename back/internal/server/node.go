package server

import (
	"context"
	"errors"
	"fmt"
	"languago/internal/pkg/closer"
	"languago/internal/pkg/config"
	errors2 "languago/internal/pkg/errors"
	"languago/internal/pkg/logger"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type (
	Node interface {
		ID() uuid.UUID
		Run()
		Stop(ctx context.Context)
		Healthcheck() []error
		SetConfig(cfg config.AbstractNodeConfig)
		ServiceBuilder(cfg config.AbstractConfig)
		ErrorsPresenter() errors2.ErrorsPersenter
		Log() logger.Logger
	}

	Service interface {
		StartService(e chan error, closer chan closer.CloseFunc)
		GracefulStop() error
		Ping(ctx context.Context) error
	}

	node struct {
		id              uuid.UUID
		config          config.AbstractNodeConfig
		logger          logger.Logger
		errorsPersenter errors2.ErrorsPersenter
		services        map[string]Service
		errorCh         chan error
		closer          closer.Closer
		closerCh        chan closer.CloseFunc
		errorsObserver  errors2.ErrorsObserver
	}

	StopFunc func(n Node) error

	Services map[string]Service

	NewNodeParams struct {
		Services        Services
		StopFuncs       []StopFunc
		Logger          logger.Logger
		Closer          closer.Closer
		ErrorsPresenter errors2.ErrorsPersenter
	}
)

func (n *node) ServiceBuilder(cfg config.AbstractConfig) {
	// just for now
	svc := map[string]Service{
		"flashcard_service": NewService(cfg),
	}
	for _, s := range svc {
		n.closer.Add(func() error {
			return s.GracefulStop()
		})
	}

	n.services = svc
}

func NewNode(args *NewNodeParams) Node {
	if args == nil {
		panic("error NewNodeParams are required.")
	}

	node := &node{
		id:              uuid.New(),
		logger:          args.Logger,
		services:        args.Services,
		errorsPersenter: args.ErrorsPresenter,
		closer:          args.Closer,
		errorsObserver:  errors2.NewErrorObserver(args.Logger),
	}

	node.LogErrors()

	return node
}

func (n *node) Run() {
	n.logger.Info("starting the node services",
		logger.LogFields{
			"node_id":       n.id.String(),
			"node_services": n.services,
		},
	)

	for name, s := range n.services {
		log.Println(name, " starting")
		s.StartService(n.errorCh, n.closerCh)
		n.logger.Info(fmt.Sprintf("service %s successfully started", name), nil)
	}
}

func (n *node) Stop(ctx context.Context) {
	n.closer.Close(ctx)
}

func (n *node) Healthcheck() []error {
	var errs []error
	var respCh chan error = make(chan error, len(n.services))
	var wg sync.WaitGroup
	defer close(respCh)

	for name, s := range n.services {
		wg.Add(1)
		go func(name string, s Service) {
			defer wg.Done()

			ctx, c := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
			defer c()

			err := s.Ping(ctx)
			switch {
			case err == nil:
				return
			case errors.Is(err, context.DeadlineExceeded):
				respCh <- fmt.Errorf("error service %s is not responding.", name)
				return
			//case errors.Is(): TODO custom errors

			default:
				respCh <- fmt.Errorf("error service %s. unknown error", name)
				return
			}
		}(name, s)
	}

	wg.Wait()
	if len(respCh) > 0 {
		for e := range respCh {
			errs = append(errs, e)
		}
	}
	return errs
}

func (n *node) ID() uuid.UUID { return n.id }

func (n *node) SetConfig(cfg config.AbstractNodeConfig) {
	n.config = cfg
}

func (n *node) Log() logger.Logger { return n.logger }

func (n *node) ErrorsPresenter() errors2.ErrorsPersenter {
	return n.errorsPersenter
}

func (n *node) LogErrors() {
	n.errorsObserver.WatchErrors(n)
}

func (n *node) ErrorChannel() chan error {
	return n.errorCh
}

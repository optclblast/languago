package server

import (
	"context"
	"errors"
	"fmt"
	"languago/internal/pkg/closer"
	"languago/internal/pkg/config"
	"languago/internal/pkg/languagoerr"
	"languago/internal/pkg/logger"
	"sync"
	"time"

	"github.com/google/uuid"
)

type (
	Node interface {
		Run()
		Stop()
		Healthcheck() []error
		SetConfig(cfg config.AbstractNodeConfig)
	}

	Service interface {
		StartService(e chan error, closer chan closer.CloseFunc)
		StopService() error
		Ping(ctx context.Context) error
	}

	node struct {
		Id       uuid.UUID
		config   config.AbstractNodeConfig
		Logger   logger.Logger
		Services map[string]Service
		ErrorCh  chan error
		StopCh   chan struct{}
		closer   closer.Closer
		closerCh chan closer.CloseFunc
	}

	StopFunc func(n Node) error

	NewNodeParams struct {
		Services  map[string]Service
		StopFuncs []StopFunc
		Logger    logger.Logger
		Closer    closer.Closer
	}
)

func NewNode(args *NewNodeParams) (Node, error) {
	if args == nil {
		return nil, fmt.Errorf("error NewNodeParams are required.")
	}
	if len(args.Services) == 0 {
		return nil, fmt.Errorf("error no services provided.")
	}

	return &node{
		Id:       uuid.New(),
		Logger:   args.Logger,
		Services: args.Services,
		closer:   args.Closer,
		closerCh: make(chan closer.CloseFunc),
	}, nil
}

func (n *node) Run() {
	go n.errorHandler()
	for name, s := range n.Services {
		s.StartService(n.ErrorCh, n.closerCh)
		n.Logger.Info(fmt.Sprintf("service %s successfully started", name), nil)
	}
	<-n.StopCh
}

func (n *node) Stop() {
	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(10*time.Second),
	)
	defer cancel()
	n.closer.Close(ctx)
}

func (n *node) Healthcheck() []error {
	var errs []error
	var respCh chan error = make(chan error, len(n.Services))
	var wg sync.WaitGroup
	defer close(respCh)

	for name, s := range n.Services {
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
			//case errors.Is(): TODO custom errors

			default:
				respCh <- fmt.Errorf("error service %s. unknown error", name)
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

func (n *node) SetConfig(cfg config.AbstractNodeConfig) {
	n.config = cfg
}

func (n *node) errorHandler() {
	for {
		err := <-n.ErrorCh
		if _, ok := err.(*languagoerr.FatalErr); ok {
			n.Logger.Warn("fatal service error: ", logger.LogFieldPair(logger.ErrorField, <-n.ErrorCh))
			// TODO restart service handler, like this -> n.RestartService() or sthg
			continue
		}
		n.Logger.Info("error: ", logger.LogFieldPair(logger.ErrorField, <-n.ErrorCh))
	}
}

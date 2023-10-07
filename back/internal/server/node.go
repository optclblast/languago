package server

import (
	"context"
	"errors"
	"fmt"
	"languago/internal/pkg/closer"
	"languago/internal/pkg/config"
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
		Log() logger.Logger
	}

	Service interface {
		StartService(e chan error, closer chan closer.CloseFunc)
		GracefulStop() error
		Ping(ctx context.Context) error
	}

	node struct {
		Id       uuid.UUID
		config   config.AbstractNodeConfig
		Logger   logger.Logger
		Services map[string]Service
		ErrorCh  chan error
		closer   closer.Closer
		closerCh chan closer.CloseFunc
	}

	StopFunc func(n Node) error

	Services map[string]Service

	NewNodeParams struct {
		Services  Services
		StopFuncs []StopFunc
		Logger    logger.Logger
		Closer    closer.Closer
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

	n.Services = svc
}

func NewNode(args *NewNodeParams) (Node, error) {
	if args == nil {
		panic("error NewNodeParams are required.")
	}

	return &node{
		Id:       uuid.New(),
		Logger:   args.Logger,
		Services: args.Services,
		closer:   args.Closer,
	}, nil
}

func (n *node) Run() {
	n.Logger.Info("starting the node services",
		logger.LogFields{
			"node_id":       n.Id.String(),
			"node_services": n.Services,
		},
	)

	for name, s := range n.Services {
		log.Println(name, " starting")
		s.StartService(n.ErrorCh, n.closerCh)
		n.Logger.Info(fmt.Sprintf("service %s successfully started", name), nil)
	}
}

func (n *node) Stop(ctx context.Context) {
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

func (n *node) ID() uuid.UUID { return n.Id }

func (n *node) SetConfig(cfg config.AbstractNodeConfig) {
	n.config = cfg
}

func (n *node) Log() logger.Logger { return n.Logger }

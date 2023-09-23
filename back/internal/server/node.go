package server

import (
	"context"
	"fmt"
	"languago/internal/pkg/closer"
	"languago/internal/pkg/config"
	"languago/internal/pkg/logger"
	"time"

	"github.com/google/uuid"
)

type (
	Node interface {
		Run()
		Stop()
		TODO()
		Healthcheck() []error
		SetConfig(cfg config.AbstractNodeConfig)
	}

	Service interface {
		StartService(e chan error, closer chan closer.CloseFunc)
		StopService() error
		Ping() bool
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
		n.Logger.Info(fmt.Sprintf("service %s successfully started", name))
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
	n.TODO()
}

func (n *node) TODO() {}

func (n *node) Healthcheck() []error {
	var errs []error
	for name, s := range n.Services {
		if ok := s.Ping(); !ok {
			errs = append(errs, fmt.Errorf("error service %s is not ok.", name))
		}
	}
	return errs
}

func (n *node) SetConfig(cfg config.AbstractNodeConfig) {
	n.config = cfg
}

func (n *node) errorHandler() {
	for {
		n.Logger.Warn("error: ", <-n.ErrorCh)
	}
}

package server

import (
	"fmt"
	"languago/internal/pkg/config"
	"languago/internal/pkg/logger"

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
		StartService(e chan error)
		StopService() error
		Ping() bool
	}

	node struct {
		Id        uuid.UUID
		config    config.AbstractNodeConfig
		Logger    logger.Logger
		Services  map[string]Service
		ErrorCh   chan error
		StopCh    chan struct{}
		StopFuncs []StopFunc
	}

	StopFunc func(n Node) error

	NewNodeParams struct {
		Services  map[string]Service
		StopFuncs []StopFunc
		Logger    logger.Logger
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
		Id:        uuid.New(),
		Logger:    args.Logger,
		Services:  args.Services,
		StopFuncs: args.StopFuncs,
	}, nil
}

func (n *node) Run() {
	go n.errorHandler()
	for name, s := range n.Services {
		s.StartService(n.ErrorCh)
		n.Logger.Info(fmt.Sprintf("service %s successfully started", name))
	}
	<-n.StopCh
}

func (n *node) Stop() {
	for _, fn := range n.StopFuncs {
		if err := fn(n); err != nil {
			n.Logger.Warn(fmt.Sprintf("error applying stop func: %s", err.Error()))
		}
	}
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

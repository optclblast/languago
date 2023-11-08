package server

import (
	"context"
	"languago/internal/pkg/config"
	errors2 "languago/internal/pkg/errors"
	"languago/internal/pkg/logger"

	"github.com/google/uuid"
)

type (
	Node interface {
		ID() uuid.UUID
		Run()
		Stop(ctx context.Context)
		SetConfig(cfg config.AbstractNodeConfig)
		ErrorsPresenter() errors2.ErrorsPersenter
		Log() logger.Logger
	}

	Service interface {
		Start(e chan error)
	}

	node struct {
		id              uuid.UUID
		config          config.AbstractNodeConfig
		logger          logger.Logger
		errorsPersenter errors2.ErrorsPersenter
		services        Services
		errorCh         chan error
		//closer          closer.Closer
		//closerCh        chan closer.CloseFunc
		errorsObserver errors2.ErrorsObserver
	}

	StopFunc func(n Node) error
	Services []Service

	NewNodeParams struct {
		//StopFuncs       []StopFunc
		Logger logger.Logger
		Config config.AbstractConfig
		//Closer          closer.Closer
		ErrorsPresenter errors2.ErrorsPersenter
	}
)

func NewNode(args *NewNodeParams) Node {
	if args == nil {
		panic("error NewNodeParams are required.")
	}

	var services Services = make(Services, 0)
	for _, serviceCfg := range args.Config.GetNodeConfig().GetServicesCfg() {
		service := NewService(args.Config, serviceCfg.GetHTTPAddress())
		services = append(services, service)
	}

	nodeId := uuid.New()

	node := &node{
		id:              nodeId,
		logger:          args.Logger,
		errorsPersenter: args.ErrorsPresenter,
		services:        services,
		//closer:          args.Closer,
		errorsObserver: errors2.NewErrorObserver(args.Logger),
	}

	errObserver := errors2.NewErrorObserver(args.Logger)
	errObserver.WatchErrors(node)

	node.LogErrors()

	return node
}

func (n *node) Run() {
	n.logger.Info("starting the node: ", logger.LogFields{
		"node_id": n.ID(),
	})

	for _, s := range n.services {
		s.Start(n.errorCh)
	}
}

func (n *node) Stop(ctx context.Context) {
	// todo
}

func (n *node) ID() uuid.UUID { return n.id }

func (n *node) SetConfig(cfg config.AbstractNodeConfig) { n.config = cfg }

func (n *node) Log() logger.Logger { return n.logger }

func (n *node) ErrorsPresenter() errors2.ErrorsPersenter { return n.errorsPersenter }

func (n *node) LogErrors() { n.errorsObserver.WatchErrors(n) }

func (n *node) ErrorChannel() chan error { return n.errorCh }

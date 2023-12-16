package server

import (
	"context"
	"languago/infrastructure/config"
	"languago/infrastructure/logger"
	errors2 "languago/pkg/errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type (
	Node interface {
		ID() uuid.UUID
		Run()
		Stop(ctx context.Context)
		SetConfig(cfg config.AbstractNodeConfig)
		ErrorsPresenter() errors2.ErrorsPersenter
		Log() *zerolog.Logger
	}

	Service interface {
		Start(e chan error)
	}

	node struct {
		id       uuid.UUID
		config   config.AbstractNodeConfig
		services Services
		log      *zerolog.Logger

		errorsPersenter errors2.ErrorsPersenter
		errorCh         chan error
		errorsObserver  errors2.ErrorsObserver

		// deprecated
		logger logger.Logger
	}

	StopFunc func(n Node) error
	Services []Service

	NewNodeParams struct {
		//StopFuncs       []StopFunc
		Log    zerolog.Logger
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
		log:             &args.Log,
		errorsPersenter: args.ErrorsPresenter,
		services:        services,
		//closer:          args.Closer,
		errorsObserver: errors2.NewErrorObserver(args.Log),
	}

	errObserver := errors2.NewErrorObserver(args.Log)
	errObserver.WatchErrors(node)

	node.LogErrors()

	return node
}

func (n *node) Run() {
	n.log.Info().Msgf("starting the node: node_id: %v", n.ID())

	for _, s := range n.services {
		s.Start(n.errorCh)
	}
}

func (n *node) Stop(ctx context.Context) {
	// todo
}

func (n *node) ID() uuid.UUID { return n.id }

func (n *node) SetConfig(cfg config.AbstractNodeConfig) { n.config = cfg }

func (n *node) Log() *zerolog.Logger {
	return n.log
}

func (n *node) ErrorsPresenter() errors2.ErrorsPersenter { return n.errorsPersenter }

func (n *node) LogErrors() { n.errorsObserver.WatchErrors(n) }

func (n *node) ErrorChannel() chan error { return n.errorCh }

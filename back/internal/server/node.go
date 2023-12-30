package server

import (
	"context"
	"languago/infrastructure/config"
	"languago/infrastructure/logger"
	errors2 "languago/pkg/errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type (
	Node interface {
		ID() uuid.UUID
		Run()
		Stop(ctx context.Context)
		SetConfig(cfg *config.Config)
		ErrorsPresenter() errors2.ErrorsPersenter
		Log() *logrus.Logger
		Version() string
	}

	Service interface {
		Start(e chan error)
	}

	node struct {
		id       uuid.UUID
		config   *config.Config
		services Services
		log      *logrus.Logger
		version  string

		errorsPersenter errors2.ErrorsPersenter
		errorCh         chan error
		errorsObserver  errors2.ErrorsObserver
	}

	StopFunc func(n Node) error
	Services []Service

	NewNodeParams struct {
		Version string
		Log     *logrus.Logger
		Logger  logger.Logger
		Config  *config.Config
		//Closer          closer.Closer
		ErrorsPresenter errors2.ErrorsPersenter
	}
)

func NewNode(args *NewNodeParams) Node {
	if args == nil {
		panic("error NewNodeParams are required.")
	}

	var services Services = make(Services, 0)
	service := NewService(args.Config, args.Config.Node.FlashcardAPI.Address, args.Log)
	services = append(services, service)

	nodeId := uuid.New()

	node := &node{
		id:              nodeId,
		version:         args.Version,
		log:             args.Log,
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
	n.log.Infof("starting the node: node_id: %v version: %s", n.ID(), n.Version())

	for _, s := range n.services {
		s.Start(n.errorCh)
	}
}

func (n *node) Stop(ctx context.Context) {
	// todo
}

func (n *node) Version() string {
	return n.version
}

func (n *node) ID() uuid.UUID { return n.id }

func (n *node) SetConfig(cfg *config.Config) { n.config = cfg }

func (n *node) Log() *logrus.Logger {
	return n.log
}

func (n *node) ErrorsPresenter() errors2.ErrorsPersenter { return n.errorsPersenter }

func (n *node) LogErrors() { n.errorsObserver.WatchErrors(n) }

func (n *node) ErrorChannel() chan error { return n.errorCh }

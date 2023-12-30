package server

import (
	"fmt"
	"languago/infrastructure/config"
	"languago/infrastructure/logger"
	"languago/infrastructure/repository"
	"languago/interface/api"

	//errors2 "languago/pkg/errors"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type (
	flashcardService struct {
		API     *api.API
		address string
		log     *logrus.Logger
		//errorsPresenter errors2.ErrorsPersenter
	}
)

func NewService(cfg *config.Config, address string, log *logrus.Logger) Service {
	dbInteractor, err := repository.NewDatabaseInteractor(cfg.Node.Database)
	if err != nil {
		panic("can't get database interactor! " + err.Error())
	}
	return &flashcardService{
		API:     api.NewAPI(cfg, log, dbInteractor),
		address: address,
		log:     logger.ProvideLogger(cfg),
	}
}

func (s *flashcardService) Start(e chan error) {
	s.log.Infof("Starting server at %v", s.address)
	go s.listen(e)
}

func (s *flashcardService) listen(e chan error) {
	err := http.ListenAndServe(s.address, s.API)
	if err != nil {
		e <- fmt.Errorf("error service runtime error: %w", err)
	}
}

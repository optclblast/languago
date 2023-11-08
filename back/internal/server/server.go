package server

import (
	"fmt"
	"languago/interface/api"
	"languago/internal/pkg/config"
	"languago/internal/pkg/logger"
	"languago/internal/pkg/repository"

	errors2 "languago/internal/pkg/errors"
	"net/http"

	_ "github.com/lib/pq"
)

type (
	flashcardService struct {
		API             *api.API
		address         string
		log             logger.Logger
		errorsPresenter errors2.ErrorsPersenter
	}
)

func NewService(cfg config.AbstractConfig, address string) Service {
	dbInteractor, err := repository.NewDatabaseInteractor(cfg.GetDatabaseConfig())
	if err != nil {
		panic("can't get database interactor! " + err.Error())
	}
	return &flashcardService{
		API:     api.NewAPI(cfg.GetLoggerConfig(), dbInteractor),
		address: address,
		log:     logger.ProvideLogger(cfg.GetLoggerConfig()),
	}
}

func (s *flashcardService) Start(e chan error) {
	s.log.Info("Starting server", logger.LogFields{
		"address": s.address,
	})
	go s.listen(e)
}

func (s *flashcardService) listen(e chan error) {
	err := http.ListenAndServe(s.address, s.API)
	if err != nil {
		e <- fmt.Errorf("error service runtime error: %w", err)
	}
}

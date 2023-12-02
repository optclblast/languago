package server

import (
	"fmt"
	"languago/infrastructure/config"
	"languago/infrastructure/logger"
	"languago/infrastructure/repository"
	"languago/interface/api"

	errors2 "languago/pkg/errors"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type (
	flashcardService struct {
		API             *api.API
		address         string
		log             zerolog.Logger
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
	s.log.Info().Msgf("Starting server at %v", s.address)
	go s.listen(e)
}

func (s *flashcardService) listen(e chan error) {
	err := http.ListenAndServe(s.address, s.API)
	if err != nil {
		e <- fmt.Errorf("error service runtime error: %w", err)
	}
}

package server

import (
	"fmt"
	"languago/internal/pkg/config"
	"languago/internal/pkg/logger"
	"languago/internal/pkg/repository"
	"languago/internal/server/api"
	"net/http"

	_ "github.com/lib/pq"
)

type (
	flashcardService struct {
		API    *api.API
		log    logger.Logger
		Config ServerConfigPresenter
	}

	ServerConfigPresenter interface {
		InitConfig() error
		ChangeConfig(up ConfigUpdateParams) error
		GetValue(key string) (interface{}, error)
		SetValue(key string, value interface{}) error
	}

	ConfigUpdateParams map[string]interface{}
)

func NewService(cfg config.AbstractConfig) Service {
	dbInteractor, err := repository.NewDatabaseInteractor(cfg.GetDatabaseConfig())
	if err != nil {
		panic("can't get database interactor! " + err.Error())
	}
	return &flashcardService{
		API: api.NewAPI(cfg.GetLoggerConfig(), dbInteractor),
		log: logger.ProvideLogger(cfg.GetLoggerConfig()),
	}
}

func (s *flashcardService) StartService(e chan error) {
	s.log.Info("Starting server")
	s.API.Init()
	go s.listen(e)
}

func (s *flashcardService) StopService() error {
	s.log.Warn("started flashcard service shutdown")
	// TODO safe shutdown
	return nil
}

func (s *flashcardService) Ping() bool {
	// TODO
	return true
}

func (s *flashcardService) listen(e chan error) {
	err := http.ListenAndServe("localhost:3300", s.API)
	if err != nil {
		e <- fmt.Errorf("error service runtime error: %w", err)
	}
}

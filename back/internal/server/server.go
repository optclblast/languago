package server

import (
	"context"
	"fmt"
	"languago/internal/pkg/closer"
	"languago/internal/pkg/config"
	"languago/internal/pkg/logger"
	"languago/internal/pkg/repository"
	"languago/internal/server/api"
	"net/http"
	"time"

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

func (s *flashcardService) StartService(e chan error, closer chan closer.CloseFunc) {
	s.log.Info("Starting server", nil)
	s.API.Init()
	go s.listen(e)
}

func (s *flashcardService) Ping(ctx context.Context) error {
	return s.API.Repo.Database().PingDB()
}

func (s *flashcardService) listen(e chan error) {
	err := http.ListenAndServe("localhost:3300", s.API)
	if err != nil {
		e <- fmt.Errorf("error service runtime error: %w", err)
	}
}

func (s *flashcardService) GracefulStop() error {
	s.log.Warn("closing database connection", logger.LogFields{
		"service": "flashcardService",
		"time":    time.Now(),
	})
	s.API.Repo.CloseConnection()
	s.API.Stop()
	return nil
}

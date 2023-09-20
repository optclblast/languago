package server

import (
	"fmt"
	"languago/internal/pkg/logger"
	"languago/internal/server/api"
	"net/http"
)

type (
	flashcardService struct {
		API     *api.API
		Storage struct{} //TODO
		log     logger.Logger
		Config  ServerConfigPresenter
	}

	ServerConfigPresenter interface {
		InitConfig() error
		ChangeConfig(up ConfigUpdateParams) error
		GetValue(key string) (interface{}, error)
		SetValue(key string, value interface{}) error
	}

	ConfigUpdateParams map[string]interface{}
)

func NewService() Service {
	return &flashcardService{
		API: api.NewAPI(),
		log: logger.ProvideLogger(nil),
	}
}

func (s *flashcardService) StartService() error {
	s.log.Info("Starting server")
	s.API.Init()
	err := http.ListenAndServe("localhost:3300", s.API)
	if err != nil {
		// Err??
		//s.log.Err(fmt.Sprintf("fatal serving error: %s", err.Error()))
		return fmt.Errorf("error service runtime error: %w", err)
	}
	return nil
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

package server

import (
	"languago/internal/pkg/logger"
	"languago/internal/server/api"
	"net/http"
)

type (
	Service struct {
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

func NewService() *Service {
	return &Service{
		API: api.NewAPI(),
		log: logger.ProvideLogger(nil),
	}
}

func (s *Service) Run() {
	s.log.Info("Starting server")
	s.API.Init()
	err := http.ListenAndServe("localhost:3300", s.API)
	if err != nil {
		// Err??
		//s.log.Err(fmt.Sprintf("fatal serving error: %s", err.Error()))
	}
}

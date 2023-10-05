package test

import (
	"languago/internal/pkg/closer"
	"languago/internal/pkg/config"
	"languago/internal/pkg/logger"
	"languago/internal/server"
	"testing"
	"time"
)

func TestNodeInitialization(t *testing.T) {
	// must be ok
	cfg := server.NewNodeParams{
		Services: map[string]server.Service{
			"flashcard": server.NewService(&config.Config{
				DatabaseCfg: &config.DatabaseConfig{
					DatabaseAddress: "localhost:5432",
					DatabaseDriver:  "postgres",
					DatabaseUser:    "postgres",
					DatabaseSecret:  "postgres",
				},
				NodeCfg: &config.NodeConfig{
					Logger: logger.ProvideLogger(&config.LoggerConfig{
						Logger: logger.NewLogrusWrapper(true, logger.EnvParam_LOCAL),
					}),
					HTTPAddress: "localhost",
					HTTPPort:    "8080",
				},
				LoggerCfg: &config.LoggerConfig{
					Logger: logger.NewLogrusWrapper(true, logger.EnvParam_LOCAL),
				},
			}),
		},
		Logger: logger.ProvideLogger(&config.LoggerConfig{
			Logger: logger.NewLogrusWrapper(true, logger.EnvParam_LOCAL),
		}),
		Closer: closer.NewCloser(),
	}

	node, err := server.NewNode(&cfg)
	if err != nil {
		t.Error(err)
	}

	go node.Run()
	time.Sleep(5 * time.Second)
	node.Stop()

	// must fail

	failCfg := cfg
	failCfg.Services = nil
	_, err = server.NewNode(&failCfg)
	if err == nil {
		t.Error(err)
	}

	// must fail

	failCfg2 := cfg
	failCfg2.Logger = nil
	_, err = server.NewNode(&failCfg)
	if err == nil {
		t.Error(err)
	}

	// must fail

	failCfg3 := cfg
	failCfg3.Closer = nil
	_, err = server.NewNode(&failCfg)
	if err == nil {
		t.Error(err)
	}
}

func TestNodeHealthcheck(t *testing.T) {
	cfg := server.NewNodeParams{
		Services: map[string]server.Service{
			"flashcard": server.NewService(&config.Config{
				DatabaseCfg: &config.DatabaseConfig{
					DatabaseAddress: "localhost:5432",
					DatabaseDriver:  "postgres",
					DatabaseUser:    "postgres",
					DatabaseSecret:  "postgres",
				},
				NodeCfg: &config.NodeConfig{
					Logger: logger.ProvideLogger(&config.LoggerConfig{
						Logger: logger.NewLogrusWrapper(true, logger.EnvParam_LOCAL),
					}),
					HTTPAddress: "localhost",
					HTTPPort:    "8080",
				},
				LoggerCfg: &config.LoggerConfig{
					Logger: logger.NewLogrusWrapper(true, logger.EnvParam_LOCAL),
				},
			}),
		},
		Logger: logger.ProvideLogger(&config.LoggerConfig{
			Logger: logger.NewLogrusWrapper(true, logger.EnvParam_LOCAL),
		}),
		Closer: closer.NewCloser(),
	}

	node, err := server.NewNode(&cfg)
	if err != nil {
		t.Error(err)
	}

	go node.Run()
	time.Sleep(5 * time.Second)

	errs := node.Healthcheck()
	if len(errs) > 0 {
		t.Error(errs)
	}
	t.Log("services are ok!")
}

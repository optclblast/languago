package config

import (
	"encoding/json"
	"fmt"
	"languago/internal/pkg/logger"
	"languago/internal/pkg/mock"
	"languago/internal/pkg/repository"
	"os"
)

var (
	CONFIG_DIR string = "./general.json"
)

type (
	cfgFileStruct struct {
		Node struct {
			HTTPAddress string `json:"http_address"`
			HTTPPort    string `json:"http_port"`
			RPCAddress  string `json:"rpc_address"`
			RPCPort     string `json:"rpc_port"`
		} `json:"node"`
		Database struct {
			DatabaseAddress string `json:"db_address"`
			DatabaseDriver  string `json:"db_driver"`
			DatabaseUser    string `json:"db_user"`
			DatabaseSecret  string `json:"db_secret"`
		} `json:"database"`
		Logger struct {
			Logger    string `json:"logger"`
			DebugMode bool   `json:"debug"`
			// LogrusParams struct {
			// 	// TODO
			// } `json:"logrus_params,omitempty"`
			// SlogParams struct {
			// } `json:"slog_params,omitempty"`
		} `json:"logger"`
	}

	AbstractConfig interface {
		GetDatabaseConfig() AbstractDatabaseConfig
		GetNodeConfig() AbstractNodeConfig
		GetLoggerConfig() AbstractLoggerConfig
	}

	AbstractDatabaseConfig interface {
		GetCredentials() repository.DBCredentials
	}

	AbstractNodeConfig interface {
		GetHTTPAddress() string
		GetRPCAddress() string
		SetLogger(l Logger)
	}

	AbstractLoggerConfig interface {
		GetLogger() Logger
	}

	Config struct {
		DatabaseCfg *DatabaseConfig
		NodeCfg     *NodeConfig
		LoggerCfg   *LoggerConfig
	}

	DatabaseConfig struct {
		DatabaseAddress string
		DatabaseDriver  string
		DatabaseUser    string
		DatabaseSecret  string
	}

	NodeConfig struct {
		log         Logger
		HTTPAddress string
		HTTPPort    string
		RPCAddress  string
		RPCPort     string
	}

	LoggerConfig struct {
		Logger Logger
	}
)

func NewConfig() AbstractConfig {
	rawCfg := &cfgFileStruct{}
	var cfg Config

	data, err := os.ReadFile(CONFIG_DIR)
	if err != nil {
		panic("error reading config file: " + err.Error())
	}

	fmt.Println(string(data))

	err = json.Unmarshal(data, rawCfg)
	if err != nil {
		panic("error unmarshaling config file: " + err.Error())
	}

	switch rawCfg.Logger.Logger {
	case "logrus":
		cfg.LoggerCfg.Logger = logger.NewLogrusWrapper(rawCfg.Logger.DebugMode)
	case "slog":
		// TODO
		//cfg.LoggerCfg.Logger = logger.NewSLogLogger(rawCfg.Logger.DebugMode)
		mock.ImplementMePanic()
	default:
		cfg.LoggerCfg.Logger = logger.NewDefaultLogger(rawCfg.Logger.DebugMode)
	}

	cfg.NodeCfg = &NodeConfig{
		log:         cfg.LoggerCfg.Logger,
		HTTPAddress: rawCfg.Node.HTTPAddress,
		HTTPPort:    rawCfg.Node.HTTPPort,
		RPCAddress:  rawCfg.Node.RPCAddress,
		RPCPort:     rawCfg.Node.RPCPort,
	}
	cfg.DatabaseCfg = (*DatabaseConfig)(&rawCfg.Database)

	return &cfg
}

func (c *Config) GetDatabaseConfig() AbstractDatabaseConfig {
	return c.DatabaseCfg
}

func (c *Config) GetNodeConfig() AbstractNodeConfig {
	return c.NodeCfg
}

func (c *Config) GetLoggerConfig() AbstractLoggerConfig {
	return c.LoggerCfg
}

func (c *DatabaseConfig) GetCredentials() repository.DBCredentials {
	return &repository.DBCred{
		DbAddress: c.DatabaseAddress,
		Driver:    c.DatabaseDriver,
		User:      c.DatabaseUser,
		Secret:    c.DatabaseSecret,
	}
}

func (c *NodeConfig) GetHTTPAddress() string {
	mock.ImplementMePanic()
	return "localhost"
}

func (c *NodeConfig) GetRPCAddress() string {
	mock.ImplementMePanic()
	return "0.0.0.0"
}

func (c *NodeConfig) SetLogger(l Logger) {
	c.log = l
}

func (c *LoggerConfig) GetLogger() Logger {
	return c.Logger
}

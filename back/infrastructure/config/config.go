package config

import (
	"bytes"
	"fmt"
	"languago/infrastructure/logger"
	"languago/infrastructure/logger/wrappers"
	"languago/infrastructure/repository"
	"log"
	"os"

	"github.com/spf13/viper"
)

type (
	AbstractConfig interface {
		GetDatabaseConfig() AbstractDatabaseConfig
		GetNodeConfig() AbstractNodeConfig
		GetLoggerConfig() AbstractLoggerConfig
	}

	AbstractDatabaseConfig interface {
		GetCredentials() repository.DBCredentials
		IsMock() bool
	}

	AbstractNodeConfig interface {
		SetLogger(l logger.Logger)
		GetServicesCfg() []AbstractServiceConfig
	}

	AbstractLoggerConfig interface {
		GetLogger() logger.Logger
	}

	AbstractServiceConfig interface {
		ServiceName() string
		GetHTTPAddress() string
	}

	Config struct {
		DatabaseCfg *DatabaseConfig
		NodeCfg     *NodeConfig
		LoggerCfg   *LoggerConfig
	}

	DatabaseConfig struct {
		isMock          bool
		DatabaseAddress string
		DatabaseDriver  string
		DatabaseUser    string
		DatabaseSecret  string
	}

	NodeConfig struct {
		Logger   logger.Logger
		Services []AbstractServiceConfig
	}

	ServiceConfig struct {
		Name    string
		Address string
	}

	LoggerConfig struct {
		Logger logger.Logger
	}
)

func InitialConfiguration() AbstractConfig {
	var config Config = Config{
		DatabaseCfg: new(DatabaseConfig),
		NodeCfg:     new(NodeConfig),
		LoggerCfg:   new(LoggerConfig),
	}
	CONFIG_DIR := os.Getenv("LANGUAGO_CONFIG_DIR")
	var CONFIG_FILE string = "general.yaml"
	if CONFIG_DIR == "" {
		log.Println("LANGUAGO_CONFIG_DIR not provided, trying to use default configuration directory")
		CONFIG_DIR = "./cfg/"
		CONFIG_FILE = "default.yaml"
	}

	viper.SetConfigType("yaml")
	yamlCfg, err := os.ReadFile(CONFIG_DIR + CONFIG_FILE)
	if err != nil {
		panic("error reading config file: " + err.Error())
	}

	err = viper.ReadConfig(bytes.NewBuffer(yamlCfg))
	if err != nil {
		panic("error reading config file with viper: " + err.Error())
	}

	dbRaw := viper.GetStringMapString("database")
	config.DatabaseCfg.DatabaseAddress = dbRaw["db_address"]
	config.DatabaseCfg.DatabaseDriver = dbRaw["db_driver"]
	config.DatabaseCfg.DatabaseUser = dbRaw["db_user"]
	config.DatabaseCfg.DatabaseSecret = dbRaw["db_secret"]
	config.DatabaseCfg.isMock = viper.GetBool("database.is_mock")

	nodeRaw := viper.GetStringMap("node.services")
	config.NodeCfg.Services = make([]AbstractServiceConfig, 0)

	for name, serviceRaw := range nodeRaw {
		var ok bool
		service := make(map[string]interface{}, len(nodeRaw))

		if service, ok = serviceRaw.(map[string]interface{}); !ok {
			panic("error init node config")
		}

		serviceConfig := new(ServiceConfig)
		var addr, port string

		for key, value := range service {
			if serviceData, ok := value.(string); ok {
				if key == "address" {
					addr = serviceData
				} else if key == "port" {
					port = serviceData
				}

				if addr != "" && port != "" {
					serviceConfig.Name = name
					serviceConfig.Address = fmt.Sprintf("%s:%s", addr, port)
					config.NodeCfg.Services = append(config.NodeCfg.Services, serviceConfig)
				}
			}
		}

		if len(config.NodeCfg.Services) == 0 {
			panic("error can't init services config")
		}
	}

	logRaw := viper.GetStringMapString("logger")
	var envValue wrappers.EnvParam = wrappers.MustToEnvParam(
		viper.GetString("logger.env"),
	)

	dbg := viper.GetBool("logger.debug")

	switch logRaw["logger"] {
	case "logrus":
		config.LoggerCfg.Logger = wrappers.NewLogrusWrapper(dbg, envValue)
	case "std":
		//config.LoggerCfg.Logger = wrappers.NewDefaultLogger(viper.GetBool("logger.debug"))
	default:
		config.LoggerCfg.Logger = wrappers.NewZerologWrapper(viper.GetBool("logger.debug"), envValue)
	}

	return &config
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

func (c *DatabaseConfig) IsMock() bool {
	return c.isMock
}

func (c *NodeConfig) GetServicesCfg() []AbstractServiceConfig {
	return c.Services
}

func (c *NodeConfig) SetLogger(l logger.Logger) {
	c.Logger = l
}

func (c *LoggerConfig) GetLogger() logger.Logger {
	return c.Logger
}

func (c *ServiceConfig) ServiceName() string {
	return c.Name
}

func (c *ServiceConfig) GetHTTPAddress() string {
	return c.Address
}

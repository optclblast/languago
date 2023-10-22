package config

import (
	"bytes"
	"fmt"
	"languago/internal/pkg/logger"
	"languago/internal/pkg/repository"
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
	}

	AbstractNodeConfig interface {
		GetHTTPAddress() string
		GetRPCAddress() string
		SetLogger(l logger.Logger)
	}

	AbstractLoggerConfig interface {
		GetLogger() logger.Logger
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
		Logger      logger.Logger
		HTTPAddress string
		HTTPPort    string
		RPCAddress  string
		RPCPort     string
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

	nodeRaw := viper.GetStringMapString("node")
	config.NodeCfg.HTTPAddress = nodeRaw["http_address"]
	config.NodeCfg.HTTPPort = nodeRaw["http_port"]
	config.NodeCfg.RPCAddress = nodeRaw["rpc_address"]
	config.NodeCfg.RPCPort = nodeRaw["rpc_port"]

	logRaw := viper.GetStringMapString("logger")
	var envValue logger.EnvParam = logger.MustToEnvParam(viper.GetString("logger.env"))

	switch logRaw["logger"] {
	case "logrus":
		config.LoggerCfg.Logger = logger.NewLogrusWrapper(viper.GetBool("logger.debug"), envValue)
	case "zap":
		config.LoggerCfg.Logger = logger.NewZapWrapper(viper.GetBool("logger.debug"), envValue)
	case "std":
		config.LoggerCfg.Logger = logger.NewDefaultLogger(viper.GetBool("logger.debug"))
	default:
		config.LoggerCfg.Logger = logger.NewZapWrapper(viper.GetBool("logger.debug"), envValue)
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

func (c *NodeConfig) GetHTTPAddress() string {
	return fmt.Sprintf("%s:%s", c.HTTPAddress, c.HTTPPort)
}

func (c *NodeConfig) GetRPCAddress() string {
	return fmt.Sprintf("%s:%s", c.RPCAddress, c.RPCPort)
}

func (c *NodeConfig) SetLogger(l logger.Logger) {
	c.Logger = l
}

func (c *LoggerConfig) GetLogger() logger.Logger {
	return c.Logger
}

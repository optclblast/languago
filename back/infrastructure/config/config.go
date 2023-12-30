package config

import (
	"bytes"
	"log"
	"os"

	"github.com/spf13/viper"
)

func InitialConfiguration() *Config {
	var config *Config = &Config{
		Node: new(NodeConfig),
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

	config.Node.Database = new(DatabaseConfig)

	dbRaw := viper.GetStringMapString("database")
	config.Node.Database.DatabaseAddress = dbRaw["db_address"]
	config.Node.Database.DatabaseDriver = dbRaw["db_driver"]
	config.Node.Database.DatabaseUser = dbRaw["db_user"]
	config.Node.Database.DatabaseSecret = dbRaw["db_secret"]
	config.Node.Database.isMock = viper.GetBool("database.is_mock")

	config.Node.Version = viper.GetString("version")

	return config
}

func defaultConfig() *Config {
	return &Config{
		Node: &NodeConfig{
			Version:    "DEBUG",
			Name:       "languago-node-dev",
			Datacenter: "dc1",
			FlashcardAPI: &FlashcardServiceConfig{
				Version:  "DEBUG",
				Address:  "localhost:3301",
				LogLevel: 6,
			},
			Database: &DatabaseConfig{
				DatabaseAddress: "localhost:5432",
				DatabaseDriver:  "postgres",
				DatabaseUser:    "languago",
				DatabaseSecret:  "languago",
				DatabaseName:    "languago",
			},
		},
	}
}

package config

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	capi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Node *NodeConfig
}

type NodeConfig struct {
	Version      string                  `json:"version`
	Name         string                  `json:"name"`
	Datacenter   string                  `json:"datacenter"`
	FlashcardAPI *FlashcardServiceConfig `json:"flashcard_api"`
	Database     *DatabaseConfig         `json:"database"`
}

type FlashcardServiceConfig struct {
	Version  string `json:"version"`
	Address  string `json:"address"`
	LogLevel int64  `json:"log_level"`
}

type DatabaseConfig struct {
	isMock          bool
	DatabaseAddress string `json:"db_address"`
	DatabaseDriver  string `json:"db_driver"`
	DatabaseUser    string `json:"db_user"`
	DatabaseSecret  string `json:"db_secret"`
	DatabaseName    string `json:"db_name"`
}

const (
	defaultConsulAddress = "127.0.0.1:8500"
	defaultDatacenter    = "dc1"
	ttl                  = time.Second * 10
)

type RegisterNodeParams struct {
	ID                string
	Name              string
	Tags              []string
	Port              int
	Address           string
	SocketPath        string
	TaggedAddresses   map[string]string
	EnableTagOverride bool
	Meta              map[string]string
	Namespace         string
	Partition         string
}

type ConsulManager struct {
	consulAddress string
	consulClient  *capi.Client
	checkID       string
	log           *logrus.Logger
}

func (c *ConsulManager) RegisterNode(params *RegisterNodeParams) error {
	c.checkID = uuid.NewString()

	check := &capi.AgentServiceCheck{
		DeregisterCriticalServiceAfter: ttl.String(),
		TLSSkipVerify:                  true,
		TTL:                            ttl.String(),
		CheckID:                        c.checkID,
	}

	reg := &capi.AgentServiceRegistration{
		ID:      params.ID,
		Name:    params.Name,
		Tags:    params.Tags,
		Port:    params.Port,
		Address: params.Address,
		Check:   check,
	}

	if err := c.consulClient.Agent().ServiceRegister(reg); err != nil {
		return fmt.Errorf("error register node. %w", err)
	}

	go func() {
		ticker := time.NewTicker(8 * time.Second)
		for {
			err := c.consulClient.Agent().UpdateTTL(c.checkID, "alive", "pass")
			if err != nil {
				c.log.Error(err)
			}
			<-ticker.C
		}
	}()

	return nil
}

func NewConsulManager(consulAddress string, checkID string, log *logrus.Logger) (*ConsulManager, error) {
	client, err := capi.NewClient(capi.DefaultConfig())
	if err != nil {
		return nil, fmt.Errorf("error create consul client. %w", err)
	}

	return &ConsulManager{
		consulAddress: consulAddress,
		checkID:       checkID,
		log:           log,
		consulClient:  client,
	}, nil
}

func (c *ConsulManager) DefaultConfiguration() {
	defaultConfig := defaultConfig()

	data, err := json.Marshal(&defaultConfig)
	if err != nil {
		c.log.Error("error set default config.", err)
		return
	}

	c.consulClient.KV().Acquire(
		&capi.KVPair{
			Key:   "config",
			Value: data,
		},
		nil,
	)
}

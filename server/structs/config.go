package structs

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type ServiceConfig struct {
	FriendlyURL    string                 `json:"services"`
	Protocol       string                 `json:"protocol"`
	ProtocolConfig map[string]interface{} `json:"protocol_config"`
}

type Config struct {
	Services map[string]ServiceConfig `json:"services"`
}

func (c *Config) ForService(name string) *ServiceConfig {
	if sc, ok := c.Services[name]; ok {
		return &sc
	}
	return nil
}

func ParseConfig(filename string) (*Config, error) {

	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	c := Config{}

	if err := json.Unmarshal(jsonBytes, &c); err != nil {
		return nil, err
	}

	return &c, nil

}

package structs

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// ServiceConfig is the configuration for a specific service
type ServiceConfig struct {
	FriendlyURL    string                 `json:"friendly_url"`
	Protocol       string                 `json:"protocol"`
	ProtocolConfig map[string]interface{} `json:"protocol_config"`
}

// Config mirrors the config.json file which holds
// a dictionary of services with their ServiceConfigs
type Config struct {
	Services map[string]ServiceConfig `json:"services"`
}

// ForService returns the ServiceConfig for a specific
// Service by name
func (c *Config) ForService(name string) *ServiceConfig {
	if sc, ok := c.Services[name]; ok {
		return &sc
	}
	return nil
}

// ParseConfig unmarshalls the file provided to the
// Config struct and returns a pointer to it
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

package config

import (
	"encoding/json"
	"io/ioutil"
)

var (
	// Endpoints holds the Endpoints that will be tested
	Endpoints EndpointList
)

// EndpointList is a map of Endpoint definitions
type EndpointList map[string]EndpointConfig

// EndpointConfig holds the configuration for a specific endpoint
type EndpointConfig struct {
	FriedlyURL string `json:"friendly_url"`
	TSConfig tsConfig `json:"teamspeak,omitempty"`
	HTTPConfig httpConfig `json:"http,omitempty"`
	MinecraftConfig minecraftConfig `json:"minecraft,omitempty"`
}

// Protocol returns the protocol used in the EndpointConfig
func (ec *EndpointConfig) Protocol() string {
	if ec.TSConfig != (tsConfig{}) {
		return "teamspeak"
	}
	if ec.HTTPConfig != (httpConfig{}) {
		return "http"
	}
	if ec.MinecraftConfig != (minecraftConfig{}) {
		return "minecraft"
	}
	return ""
}

type tsConfig struct {
	QueryURL string `json:"query_url"`
}

type httpConfig struct {
	SuccessCodes string `json:"success_codes"`
	TestURL string `json:"test_url,omitempty"`
	Credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"credentials,omitempty"`
}

type minecraftConfig struct {
	URL string `json:"url"`
}

// LoadEndpointsFromFile returns a Config type initialized from the json file
func LoadEndpointsFromFile(path string) (EndpointList, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var el EndpointList
	if err = json.Unmarshal(file, &el); err != nil {
		return nil, err
	}
	return el, nil
}

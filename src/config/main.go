package config

import (
	"encoding/json"
	"io/ioutil"
)

// Endpoints holds the Endpoints that will be tested
var Endpoints EndpointList

// EndpointList is a map of Endpoint definitions.
// The string key is the name of the service
type EndpointList map[string]EndpointConfig

// EndpointConfig holds the friendly URL (the one listed in the API response)
// and additional protocol specific configuration (like credentials)
type EndpointConfig struct {
	FriedlyURL      string           `json:"friendly_url"`
	HTTPConfig      *httpConfig      `json:"http,omitempty"`
	TSConfig        *tsConfig        `json:"teamspeak,omitempty"`
	MinecraftConfig *minecraftConfig `json:"minecraft,omitempty"`
}

// Protocol returns the protocol used in the EndpointConfig
func (ec *EndpointConfig) Protocol() string {
	if ec.HTTPConfig != nil {
		return "http"
	}
	if ec.TSConfig != nil {
		return "teamspeak"
	}
	if ec.MinecraftConfig != nil {
		return "minecraft"
	}
	return ""
}

type httpConfig struct {
	SuccessCodes string `json:"success_codes"`
	// If the test URL is empty, the friendly URL will be used
	TestURL     string `json:"test_url,omitempty"`
	// If the Credentials are set, basicauth will be used to
	// authenticate against the webserver. This is nil by default
	// as it is a struct pointer.
	Credentials *struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"credentials,omitempty"`
}

type tsConfig struct {
	QueryURL string `json:"query_url"`
}

type minecraftConfig struct {
	URL string `json:"url"`
}

// LoadEndpointsFromFile returns a Config type initialized from the json file
func LoadEndpointsFromFile(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var el EndpointList
	if err = json.Unmarshal(file, &el); err != nil {
		return err
	}
	Endpoints = el
	return nil
}

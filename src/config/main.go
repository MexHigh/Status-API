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
	Protocol   string `json:"protocol"`
	URL        string `json:"url"`
	TSQueryURL string `json:"query_url,omitempty"`
	Credentials struct {
		Username string `json:"username"`
		Password string	`json:"password"`
	} `json:"credentials,omitempty"`
	SuccessOn  []int  `json:"success_on,omitempty"`
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

package checker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)


var (
	// Endpoints holds the Endpoints that will be tested
	Endpoints EndpointList
)

// EndpointList is a map of Endpoint definitions.
// The string key is the name of the service
type EndpointList map[string]Endpoint

// UnmarshalJSON fulfills the json.Unmarshaler interface for EndpointList.
// It therefore will be unmarshaled with the help of this function.
func (el *EndpointList) UnmarshalJSON(b []byte) error {

	*el = make(map[string]Endpoint)

	// tempEndpointList cannot be EndpointList. This would cause
	// an infinit recursion because the next call to json.Unmarshal
	// would call this function again.
	var tempEndpointList map[string]Endpoint
	if err := json.Unmarshal(b, &tempEndpointList); err != nil {
		return err
	}

	for name, endpoint := range tempEndpointList {

		var protoConfig defaultableConfig

		switch endpoint.Protocol.Type {
		case "http", "https":
			protoConfig = &HTTPConfig{}
		case "teamspeak", "ts":
			protoConfig = &TSConfig{}
		case "minecraft", "mc":
			protoConfig = &MinecraftConfig{}
		default:
			return fmt.Errorf(`Service "%s": Protocol "%s" not supported`, name, endpoint.Protocol.Type)
		}

		if err := json.Unmarshal(endpoint.Protocol.ConfigRaw, &protoConfig); err != nil {
			return err
		}
		endpoint.Protocol.Config = protoConfig
		(*el)[name] = endpoint

	}
	return nil
}

// interface guard
var _ json.Unmarshaler = (*EndpointList)(nil)

func (el *EndpointList) setDefaults() {
	for _, config := range *el {
		config.Protocol.setDefaults()
	}
}

// Endpoint holds the friendly URL (the one listed in the API response)
// and additional protocol specific configuration (like credentials)
type Endpoint struct {
	FriedlyURL string   `json:"friendly_url"`
	Status map[string]string
	Protocol   Protocol `json:"protocol"`
}

// EndpointStatus -
type EndpointStatus map[string]string

// Protocol holds information about an Endpoints protocol with which it will
// be tested. The ConfigRaw field is used to unmarshal the config.json and should
// not be used. Use the Config field (must be type asserted) instead.
type Protocol struct {
	Type      string          `json:"type"`
	ConfigRaw json.RawMessage `json:"config"`
	Config    defaultableConfig
}

func (p Protocol) setDefaults() {
	p.Config.setDefaults()
}

type defaultableConfig interface{
	setDefaults() 
}

// TSConfig is the config struct for testing if a Teamspeak 3/5 server is online
type TSConfig struct {
	QueryURL string `json:"query_url"`
}

func (t *TSConfig) setDefaults() {
	return
}

// MinecraftConfig is the config struct for testing if a minecraft server is up.
// It also displays the number of online players in the API response.
type MinecraftConfig struct {
	URL string `json:"url"`
}

func (m *MinecraftConfig) setDefaults() {
	return
}

// LoadEndpointsFromFile unmarshals the config file
func LoadEndpointsFromFile(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if !json.Valid(file) {
		return fmt.Errorf("File %s is not correctly JSON encoded", path)
	}
	var el EndpointList
	if err := json.Unmarshal(file, &el); err != nil {
		panic(err)
	}
	el.setDefaults()
	Endpoints = el
	return nil
}

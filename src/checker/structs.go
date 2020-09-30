package checker

import (
	"errors"
	"encoding/json"
)

// Endpoints holds the Endpoints that will be tested
var Endpoints EndpointList

// EndpointList is a map of Endpoint definitions.
// The string key is the name of the service
type EndpointList map[string]*Endpoint

// GetEndpoint -
func (el *EndpointList) GetEndpoint(name string) (*Endpoint, error) {
	if endpoint, ok := (*el)[name]; ok {
		return endpoint, nil
	}
	return nil, errors.New("Endpoint does not exist in config.json")
}

func (el *EndpointList) setDefaults() {
	for _, config := range *el {
		config.Protocol.setDefaults()
	}
}

// Endpoint holds the friendly URL (the one listed in the API response)
// and additional protocol specific configuration (like credentials)
type Endpoint struct {
	FriedlyURL string         `json:"friendly_url"`
	Protocol   Protocol       `json:"protocol"`
	Status     EndpointStatus `json:"-"`
}

// CheckIfUp -
func (e *Endpoint) CheckIfUp() error {
	switch p := e.Protocol.Type; p {
	case "http":
		if err := checkHTTP(e.FriedlyURL, e); err != nil {
			return err
		}
	case "teamspeak":
		if err := checkTeamspeak(e.FriedlyURL, e); err != nil {
			return err
		}
	case "minecraft":
		if err := checkMinecraft(e.FriedlyURL, e); err != nil {
			return err
		}
	default:
		panic("Error while switching through endpoint protocol type on startup")
	}
	return nil
}

// EndpointStatus -
type EndpointStatus map[string]string

// JSON -
func (es *EndpointStatus) JSON() ([]byte, error) {
	json, err := json.MarshalIndent(*es, "", "    ")
		if err != nil {
			return nil, err
		}
	return json, nil
}

// Protocol holds information about an Endpoints protocol with which it will
// be tested. The ConfigRaw field is used to unmarshal the config.json and should
// not be used. Use the Config field (must be type asserted) instead.
type Protocol struct {
	Type      string          `json:"type"`
	ConfigRaw json.RawMessage `json:"config"`
	Config    defaultableConfig `json:"-"`
}

func (p Protocol) setDefaults() {
	p.Config.setDefaults()
}

type defaultableConfig interface{
	setDefaults()
}

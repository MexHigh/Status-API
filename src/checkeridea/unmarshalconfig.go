package checkeridea

import (
	"reflect"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

// UnmarshalJSON fulfills the json.Unmarshaler interface for EndpointList.
// It therefore will be unmarshaled with the help of this function.
func (el *EndpointList) UnmarshalJSON(b []byte) error {

	*el = make(map[string]*Endpoint)

	var tempEndpointList map[string]interface{}
	if err := json.Unmarshal(b, &tempEndpointList); err != nil {
		return err
	}

	for name, endpointIF := range tempEndpointList {

		var actualEndpoint Endpoint

		if err := 

		// actualEndpoint is not a pointer to Endpoint (but must be
		// filled with pointers) because the interfaces methods are
		// implemented for the respecting pointer receivers of the
		// structs.
		/*var actualEndpoint Endpoint

		switch endpoint.Protocol {
		case "http", "https":
			actualEndpoint = &HTTPEndpoint{}
		case "teamspeak", "ts":
			actualEndpoint = &TeamspeakEndpoint{}
		case "minecraft", "mc":
			actualEndpoint = &MinecraftEndpoint{}
		default:
			return fmt.Errorf(`Service "%s": Could not match protocol`, name)
		}

		if err := json.Unmarshal(endpoint.ConfigRaw, &actualEndpoint); err != nil {
			return err
		}
		endpoint.Protocol.Config = protoConfig
		(*el)[name] = endpoint*/

	}
	return nil
}

// interface guard
var _ json.Unmarshaler = (*EndpointList)(nil)

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
	for _,  endpoint := range el {
		(*endpoint).SetDefaults()
	}
	Endpoints = el
	return nil
}
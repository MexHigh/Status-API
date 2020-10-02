package checker

import (
	"io/ioutil"
	"fmt"
	"encoding/json"
)

// UnmarshalJSON fulfills the json.Unmarshaler interface for EndpointList.
// It therefore will be unmarshaled with the help of this function.
func (el *EndpointList) UnmarshalJSON(b []byte) error {

	*el = make(map[string]*Endpoint)

	// tempEndpointList cannot be EndpointList. This would cause
	// an infinit recursion because the next call to json.Unmarshal
	// would call this function again.
	var tempEndpointList map[string]*Endpoint
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
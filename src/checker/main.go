package checker

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"status-api/config"
)

// Status holds the status information about every service registered
var Status EndpointsStatusList = make(EndpointsStatusList)

// EndpointStatus defines the status of one endpoint status
type EndpointStatus map[string]string

// JSON returns a json-formatted []byte of an EndpointStatus
func (es EndpointStatus) JSON() ([]byte, error) {
	json, err := json.MarshalIndent(es, "", "    ")
	if err != nil {
		return nil, err
	}
	return json, nil
}

// EndpointsStatusList is a map containing all services and their statuses
type EndpointsStatusList map[string]EndpointStatus

// JSON returns a json-formatted []byte of an EndpointStatusList
func (esl EndpointsStatusList) JSON() ([]byte, error) {
	json, err := json.MarshalIndent(esl, "", "    ")
	if err != nil {
		return nil, err
	}
	return json, nil
}

// GetEndpoint returns a OneEndpointStatus of a AllEndpointsStatus
func (esl EndpointsStatusList) GetEndpoint(name string) (EndpointStatus, error) {
	if os, ok := esl[name]; ok {
		return os, nil
	}
	return nil, errors.New("Endpoint does not exist in config.json")
}

// CheckService checks, if an endpoint returns one of the specified status codes
func CheckService(name string, endpoint config.EndpointConfig) error {
	switch p := endpoint.Protocol(); p {
	case "http":
		if err := checkHTTP(name, endpoint); err != nil {
			return err
		}
	case "teamspeak":
		if err := checkTeamspeak(name, endpoint); err != nil {
			return err
		}
	case "minecraft":
		if err := checkMinecraft(name, endpoint); err != nil {
			return err
		}
	default:
		return errors.New("Protocol in config file not supported")
	}
	return nil
}

// CheckAllServices checks all services mentioned in the config.json
func CheckAllServices() error {
	for name, endpoint := range config.Endpoints {
		err := CheckService(name, endpoint)
		if err != nil {
			return err
		}
	}
	return nil
}

// Updater keeps track of the records in the config.json file.
// This method is intended to be ran as goroutine (blocks until the next interval)
func Updater(interval int) {
	log.Println("Starting updater routine with an interval of " + strconv.Itoa(interval) + " seconds")
	for {
		if err := CheckAllServices(); err != nil {
			panic(err)
		}
		// wait
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

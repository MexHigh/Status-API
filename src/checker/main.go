package checker

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"status-api/config"
)

var (
	// Status holds the status information about every service registered
	Status AllEndpointsStatus = make(map[string]OneEndpointStatus)
)

// OneEndpointStatus defines the status of one endpoint status
type OneEndpointStatus map[string]string

// AllEndpointsStatus defines a map containing all services and their statuses
type AllEndpointsStatus map[string]OneEndpointStatus

// JSON returns a json-formatted []byte
func (as AllEndpointsStatus) JSON() ([]byte, error) {
	json, err := json.MarshalIndent(as, "", "    ")
	if err != nil {
		return nil, err
	}
	return json, nil
}

// JSON returns a json-formatted []byte
func (os OneEndpointStatus) JSON() ([]byte, error) {
	json, err := json.MarshalIndent(os, "", "    ")
	if err != nil {
		return nil, err
	}
	return json, nil
}

// GetEndpoint returns a OneEndpointStatus of a AllEndpointsStatus
func (as AllEndpointsStatus) GetEndpoint(name string) (OneEndpointStatus, error) {
	if os, ok := as[name]; ok {
		return os, nil
	}
	return nil, errors.New("Endpoint does not exist in config.json")
}

// CheckService checks, if an endpoint returns one of the specified status codes
func CheckService(name string, endpoint config.EndpointConfig) error {
	r, err := http.Get(endpoint.URL)
	if err != nil {
		return err
	}
	for _, statusCode := range endpoint.SuccessOn {
		if r.StatusCode == statusCode {
			Status[name] = map[string]string{
				"url":    endpoint.URL,
				"status": "up",
				"code":   strconv.Itoa(r.StatusCode),
			}
			return nil
		}
	}
	Status[name] = map[string]string{
		"url":    endpoint.URL,
		"status": "down",
		"code":   strconv.Itoa(r.StatusCode),
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
// This method is intended to be ran as goroutine (blocks )
func Updater(interval int) {
	log.Println("Starting updater routine with an interval of " + strconv.Itoa(interval) + " seconds")
	for {
		CheckAllServices()
		// wait
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

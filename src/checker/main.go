package checker

import (
	"log"
	"time"
	"errors"
	"strconv"
	"encoding/json"
	"net/http"
	"../config"
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
func (as AllEndpointsStatus) GetEndpoint(endpoint string) (OneEndpointStatus, error) {
	if os, ok := as[endpoint]; ok {
		return os, nil
	}
	return nil, errors.New("Endpoint does not exist in config.json")
}

// CheckService checks, if an endpoint returns one of the specified status codes
func CheckService(endpoint string, statusCodes []int) error {
	r, err := http.Get(endpoint)
	if err != nil {
		return err
	}
	for _, statusCode := range statusCodes {
		if r.StatusCode == statusCode {
			Status[endpoint] = map[string]string{
					"status": "up",
					"code": strconv.Itoa(r.StatusCode),
			}
			return nil
		}
	}
	Status[endpoint] = map[string]string{
			"status": "down",
			"code": strconv.Itoa(r.StatusCode),
	}
	return nil
}

// CheckAllServices checks all services mentioned in the config.json
func CheckAllServices() error {
	for endpoint, statusCodes := range config.Conf {
		err := CheckService(endpoint, statusCodes)
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
		CheckAllServices();
		// wait
		time.Sleep(time.Duration(interval) * time.Second);
	}
}
package checkeridea

import (
	"time"
	"log"
	"errors"
)

// Endpoints holds the Endpoints that will be tested
var Endpoints EndpointList

// EndpointList is a map of Endpoint definitions.
// The string key is the name of the service.
type EndpointList map[string]*Endpoint

// GetEndpoint -
func (el *EndpointList) GetEndpoint(name string) (*Endpoint, error) {
	if endpoint, ok := (*el)[name]; ok {
		return endpoint, nil
	}
	return nil, errors.New("Endpoint does not exist in config.json")
}

// CheckAll -
func (el *EndpointList) CheckAll() error {
	for _, endpoint := range *el {
		if err := (*endpoint).Check(); err != nil {
			return err
		}
	}
	return nil
}

// Updater keeps track of the records in the config.json file.
// This method is intended to be ran as goroutine (blocks until the next interval)
func Updater(interval int) {
	log.Printf("Starting updater routine with an update interval of %d seconds", interval)
	for {
		if err := Endpoints.CheckAll(); err != nil {
			panic(err)
		}
		// wait
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
package schedules

import (
	"fmt"
	"time"

	"status-api/database"
	"status-api/structs"
)

// Checkers will be filled by protocols.Register()
// called from within a checker to register themselfs
var Checkers = make(map[string]structs.Checker)

func runChecks(config *structs.Config) {

	// ResultWithName is only used here as the channel
	// type so that only one channel is necessary
	type ResultWithName struct {
		Name   string
		Result structs.CheckResult
	}

	// buffered channel for test results
	resultsChan := make(chan ResultWithName, len(config.Services))

	for name, config := range config.Services {
		// Invoke goroutines for checking every service
		go func(name string, config structs.ServiceConfig) {

			if config.ProtocolConfig == nil { // create protocolConfig if non-existent
				config.ProtocolConfig = make(map[string]interface{})
			}

			var r structs.CheckResult
			var err error

			// perform the check if
			// TODO find out what i wanted to say in the previous line
			if c, ok := Checkers[config.Protocol]; ok {
				r, err = c.Check(name, &config)
			} else {
				panic(fmt.Sprintf("Protocol %s not supported", config.Protocol))
			}
			// on unhandled error from Check method, report down with error reason
			if err != nil {
				r = structs.CheckResult{
					Status: "down",
					URL:    config.FriendlyURL,
					Reason: err.Error(),
				}
			}

			// Write the result to the channel
			resultsChan <- ResultWithName{
				Name:   name,
				Result: r,
			}

		}(name, config)
	}

	results := structs.CheckResults{
		At:       time.Now(),
		Services: make(map[string]structs.CheckResult),
	}
	// collect results from the channel
	for range config.Services {
		res := <-resultsChan
		results.Services[res.Name] = res.Result
	}
	// every service just returns one result,
	// so closing the channel here is fine
	close(resultsChan)

	model := &structs.CheckResultsModel{
		Data: results,
	}
	database.Con.Create(model)

}

package schedules

import (
	"fmt"
	"time"

	"status-api/database"
	"status-api/schedules/checkers"
	"status-api/structs"
)

func runChecks(config *structs.Config) {

	// ResultWithName is only used here as the channel
	// type so that only one channel is necessary
	type ResultWithName struct {
		Name   string
		Result structs.Result
	}

	// buffered channel for test results
	resultsChan := make(chan ResultWithName, len(config.Services))

	for name, config := range config.Services {
		// Invoke goroutines for checking every service
		go func(name string, config structs.ServiceConfig) {

			if config.ProtocolConfig == nil {
				config.ProtocolConfig = make(map[string]interface{})
			}

			var r structs.Result
			var err error

			switch proto := config.Protocol; proto {
			case "http":
				r, err = checkers.HTTP{}.Check(name, &config)
			default:
				panic(fmt.Sprintf("Protocol %s not supported", proto)) // TODO should be handled a different way
			}

			if err != nil {
				panic(err) // TODO should be handled a different way
			}

			// Idea: Add an error field to ResultWithName and check when
			// reading from the channel. If there is one, set an error
			// key in the Misc field of Result and set the status to "unknown"

			// Write the result to the channel
			resultsChan <- ResultWithName{
				Name:   name,
				Result: r,
			}

		}(name, config)
	}

	results := structs.Results{
		At:       time.Now(),
		Services: make(map[string]structs.Result),
	}
	// collect results from the channel
	for range config.Services {
		res := <-resultsChan
		results.Services[res.Name] = res.Result
	}
	// every service just returns one result,
	// so closing the channel here is fine
	close(resultsChan)

	database.Con.Create(&results)

}

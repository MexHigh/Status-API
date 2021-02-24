package schedules

import (
	"fmt"
	"sync"
	"time"

	"status-api/schedules/checkers"
	"status-api/structs"
)

func runCheck(config *structs.Config) {

	type ResultWithName struct {
		Name   string
		Result structs.Result
	}

	resultsChan := make(chan ResultWithName, len(config.Services)) // buffered channel for test results
	var wg sync.WaitGroup

	for name, config := range config.Services {
		wg.Add(1)
		go func(name string, config structs.ServiceConfig) { // check router function

			defer wg.Done()

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

			resultsChan <- ResultWithName{
				Name:   name,
				Result: r,
			}

		}(name, config)
	}
	wg.Wait() // await all goroutines (like Promise.all())
	close(resultsChan)

	results := structs.Results{}
	results.At = time.Now()
	results.Services = make(map[string]structs.Result)
	for result := range resultsChan { // read from channel until empty
		results.Services[result.Name] = result.Result
	}

	// TODO store in Database

}

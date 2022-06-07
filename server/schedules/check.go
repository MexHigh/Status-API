package schedules

import (
	"fmt"
	"log"
	"time"

	"status-api/database"
	"status-api/notify"
	"status-api/protocols"
	"status-api/structs"
)

const checkIntervalDefault = 120 // two minutes

// StartCheckTriggerJob starts the goroutine
// to trigger regular checks
func StartCheckTriggerJob(config *structs.Config) {

	// defer initial check to let the API server start
	time.Sleep(time.Duration(5) * time.Second)

	// set default interval if not defined in config.json
	checkInterval := config.CheckInterval
	if checkInterval == 0 {
		log.Println("\"check_interval\" not defined in config -> using default")
		checkInterval = checkIntervalDefault
	}
	log.Println("Check interval set to", checkInterval, "seconds")

	ran := make(chan time.Duration)

	_, err := scheduler.Every(checkInterval).Seconds().SingletonMode().Do(func() {
		start := time.Now()
		runChecks(config)
		ran <- time.Since(start)
	})
	if err != nil {
		panic(err)
	}

	for {
		log.Printf("Did checks (took %d ms)", (<-ran).Milliseconds())
	}

}

func runChecks(config *structs.Config) {

	// buffered channel for test results
	resultsChan := make(chan structs.CheckResultWithNameAndTime, len(config.Services))

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
			if c := protocols.GetChecker(config.Protocol); c != nil {
				r, err = c.Check(name, &config)
			} else {
				panic(fmt.Sprintf("Protocol %s not supported", config.Protocol))
			}
			// on unhandled error from Check method, report down with error reason
			if err != nil {
				r = structs.CheckResult{
					Status: structs.Down,
					URL:    config.FriendlyURL,
					Reason: err.Error(),
				}
			}

			resultWithName := structs.CheckResultWithNameAndTime{
				Name:   name,
				Result: r,
			}

			// Write the result to the channel
			resultsChan <- resultWithName

			// notify
			if resultWithName.Result.Status == structs.Down {
				notify.ReportDown(&resultWithName)
			} else {
				notify.ReportUp(&resultWithName)
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

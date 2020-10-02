package checker

import (
	"time"
	"log"
	"strconv"
)

func checkAllServices() error {
	for _, endpoint := range Endpoints {
		endpoint.Status = make(EndpointStatus)
		if err:= endpoint.CheckIfUp(); err != nil {
			panic(err)
		}
	}
	return nil
}

// Updater keeps track of the records in the config.json file.
// This method is intended to be ran as goroutine (blocks until the next interval)
func Updater(interval int) {
	log.Println("Starting updater routine with an update interval of " + strconv.Itoa(interval) + " seconds")
	for {
		if err := checkAllServices(); err != nil {
			panic(err)
		}
		// wait
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
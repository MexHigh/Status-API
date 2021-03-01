package schedules

import (
	"log"
	"time"

	"status-api/structs"
)

// ArchiveTriggerRoutine starts the goroutine
// to trigger archiving of all checks
func ArchiveTriggerRoutine(config *structs.Config) {

	//runArchiving(config)

	time.Sleep(time.Duration(2) * time.Second) // Defer start

	nextMidnight := func() time.Time { // at 23:59:00
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 0, 0, now.Location())
		if now.After(next) { // if between 23:59:00 and 00:00:00
			next.Add(time.Duration(24 * time.Hour))
		}
		return next
	}

	for {

		midnight := nextMidnight()
		log.Println("Next archiving scheduled at", midnight)

		time.Sleep(time.Until(midnight))

		log.Println("Archiving trigger fired. Running archiving")
		runArchiving(config)
		log.Println("Archiving done")

	}

}

const checkIntervalDefault = 120 // two minutes

// CheckTriggerRoutine starts the goroutine
// to trigger regular checks
func CheckTriggerRoutine(config *structs.Config) {

	time.Sleep(time.Duration(4) * time.Second) // Defer initial check

	log.Println("Running initial checks")
	runChecks(config)
	log.Println("Checks done")

	time.Sleep(time.Duration(1) * time.Second)

	checkInterval := config.CheckInterval
	if checkInterval == 0 {
		log.Println("\"check_interval\" not defined in config -> using default")
		checkInterval = checkIntervalDefault
	}
	log.Println("Check interval set to", checkInterval, "seconds")

	for {
		time.Sleep(time.Duration(checkInterval) * time.Second)
		log.Println("Check trigger fired. Running checks")
		runChecks(config)
		log.Println("Checks done")
	}

}

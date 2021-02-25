package schedules

import (
	"log"
	"time"

	"status-api/structs"
)

// ArchiveTriggerRoutine starts the goroutine schedule
// to trigger regular archiving
func ArchiveTriggerRoutine(config *structs.Config) {

	time.Sleep(time.Duration(2) * time.Second) // Defer start

	nextMidnight := func() time.Time {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 3, 0, 0, now.Location())
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

// CheckTriggerRoutine starts the goroutine schedule
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

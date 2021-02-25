package schedules

import (
	"log"
	"time"

	"status-api/structs"
)

// ArchiveTriggerRoutine starts the goroutine schedule
// to trigger regular archiving
func ArchiveTriggerRoutine(config *structs.Config) {
	// TODO add actual scheduling
	runArchiving(config)
}

const checkIntervalDefault = 120 // two minutes

// CheckTriggerRoutine starts the goroutine schedule
// to trigger regular checks
func CheckTriggerRoutine(config *structs.Config, checkInterval int) {

	time.Sleep(time.Duration(4) * time.Second) // Defer initial check

	log.Println("Running initial checks")
	runChecks(config)
	log.Println("Checks done")

	time.Sleep(time.Duration(1) * time.Second)

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

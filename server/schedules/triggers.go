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

// CheckTriggerRoutine starts the goroutine schedule
// to trigger regular checks
func CheckTriggerRoutine(config *structs.Config, interval int) {

	log.Println("Running initial checks")
	runChecks(config)
	log.Println("Checks done")

	for {
		time.Sleep(
			time.Duration(interval) * time.Second,
		)
		log.Println("Check trigger fired. Running checks")
		runChecks(config)
		log.Println("Checks done")
	}

}

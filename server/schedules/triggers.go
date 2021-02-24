package schedules

import (
	"status-api/structs"
)

// ArchiveTriggerRoutine starts the goroutine schedule
// to trigger regular archiving
func ArchiveTriggerRoutine(config *structs.Config) {
	// TODO add actual scheduling
	runArchive(config)
}

// CheckTriggerRoutine starts the goroutine schedule
// to trigger regular checks
func CheckTriggerRoutine(config *structs.Config, interval int) {
	// TODO add actual scheduling
	runCheck(config)
}

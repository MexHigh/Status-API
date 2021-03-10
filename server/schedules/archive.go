package schedules

import (
	"log"
	"math"
	"time"

	"status-api/database"
	"status-api/structs"
)

// StartArchiveTriggerJob starts the goroutine
// to trigger archiving of all checks at 23:59 every day
func StartArchiveTriggerJob(config *structs.Config) {

	time.Sleep(time.Duration(3) * time.Second)

	ran := make(chan time.Duration)

	job, err := scheduler.Every(1).Day().At("23:59").SingletonMode().Do(func() {
		start := time.Now()
		runArchiving(config)
		ran <- time.Since(start)
	})
	if err != nil {
		panic(err)
	}

	for {
		log.Println("Next archiving scheduled at", job.NextRun())
		log.Printf("Did archiving (took %d ms)", (<-ran).Milliseconds())
	}

}

func runArchiving(config *structs.Config) {

	dayOnly := func() time.Time {
		now := time.Now()
		return time.Date(
			now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location(),
		)
	}

	// get rows
	rows, err := database.Con.Model(&structs.CheckResultsModel{}).Rows() // TODO maybe only get todays entries?
	if err != nil {
		panic(err) // should never happen
	}

	// count up-/downtimes
	type UpsAndDowns struct {
		Ups, Downs int
		Downtimes  []structs.Downtime
	}

	udsMap := make(map[string]*UpsAndDowns) // maps service name to UpsAndDowns counter
	idsToDelete := make([]int, 0)           // holds the IDs used to create the archive result to be deleted later

	for rows.Next() {

		arch := &structs.CheckResultsModel{}
		database.Con.ScanRows(rows, arch)

		for name, service := range arch.Data.Services {
			if _, ok := udsMap[name]; !ok { // checks if upsAndDowns object exists
				udsMap[name] = &UpsAndDowns{
					Downtimes: make([]structs.Downtime, 0),
				}
			}
			if service.Status == "up" {
				udsMap[name].Ups++
			} else {
				udsMap[name].Downs++
				udsMap[name].Downtimes = append(udsMap[name].Downtimes, structs.Downtime{
					At:     arch.Data.At,
					Reason: service.Reason,
				})
			}
		}

		idsToDelete = append(idsToDelete, int(arch.ID))

	}

	rows.Close()
	if len(idsToDelete) > 0 { // suppresses error logs
		// delete all entries, that were used to compose the archive entry
		database.Con.Delete(&structs.CheckResultsModel{}, idsToDelete)
	}

	// Calulate availabilities etc.
	resultServices := make(map[string]structs.ArchiveResult)

	for name, uds := range udsMap {
		availabilityFull := float64(uds.Ups) / float64(uds.Ups+uds.Downs)
		availability := math.Round(availabilityFull*100) / 100 // rounds to two decimal places
		var status string
		if availability > 0.9 {
			status = "up"
		} else if availability > 0.7 {
			status = "problems"
		} else {
			status = "down"
		}
		resultServices[name] = structs.ArchiveResult{
			Status:       status,
			Availability: availability,
			Downtimes:    uds.Downtimes,
		}
	}

	// store in database
	resultModel := &structs.ArchiveResultsModel{
		Data: structs.ArchiveResults{
			At:       dayOnly(),
			Services: resultServices,
		},
	}
	database.Con.Create(resultModel)

	// delete entries older than 30
	var older []*structs.ArchiveResultsModel
	// SELECT * FROM archive_results_models ORDER BY id desc LIMIT 1000 OFFSET 30;
	database.Con.Model(&structs.ArchiveResultsModel{}).Order("id desc").Limit(1000).Offset(30).Scan(&older)
	if len(older) > 0 {
		// DELETE FROM archive_results_models WHERE ID IN (...)
		database.Con.Delete(&older)
	}

}

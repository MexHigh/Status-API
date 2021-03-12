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

	// counts up-/downtimes (only used locally)
	type upsAndDowns struct {
		Ups, Downs int
		Downtimes  []structs.Downtime
	}

	udsMap := make(map[string]*upsAndDowns)   // maps service name to UpsAndDowns counter
	lastDownReason := make(map[string]string) // maps the last reason of downtime to a service name
	idsToDelete := make([]int, 0)             // holds the IDs used to create the archive result to be deleted later

	for rows.Next() { // iterates over CheckResults

		currentCr := &structs.CheckResultsModel{}
		database.Con.ScanRows(rows, currentCr)

		for name, service := range currentCr.Data.Services { // iterates over one service in a CheckResults (yields service name and CheckResult)

			var uds *upsAndDowns                // local uds object (pointer to map entry) for shorter variable names
			if locUds, ok := udsMap[name]; ok { // checks if upsAndDowns object exists
				uds = locUds
			} else { // otherwise create a new one
				locUds = &upsAndDowns{
					Downtimes: make([]structs.Downtime, 0),
				}
				udsMap[name] = locUds // save new upsAndDowns in map
				uds = udsMap[name]    // reference the newly created upsAndDowns
			}

			// this is where the archiving magic happens
			if service.Status == structs.Up {
				uds.Ups++
				// reset lastDownReason to something random to prevent the concatenation of
				// downtimes if a previous error occures at a later time again
				lastDownReason[name] = "V2VyIGRhcyBsaWVzdCBpc3QgZG9vZiA7KQ=="
			} else {
				uds.Downs++
				// check if previous downtime should be ajusted, or if a new one should be created
				if lastDownReason[name] == service.Reason {
					// if the reason of the last downtime is the same as this one just
					// overwrite the latter 'To' field to the time of the current downtime
					lastDowntime := &uds.Downtimes[len(uds.Downtimes)-1] // gets pointer to last downtime
					lastDowntime.To = currentCr.Data.At
				} else {
					// if the reason is a new one, create a new downtime object and overwrite the lastDownReason
					uds.Downtimes = append(uds.Downtimes, structs.Downtime{
						From:   currentCr.Data.At,
						To:     currentCr.Data.At,
						Reason: service.Reason,
					})
					lastDownReason[name] = service.Reason
				}
			}

		}

		idsToDelete = append(idsToDelete, int(currentCr.ID))

	}

	rows.Close()
	if len(idsToDelete) > 0 { // suppresses error logs
		// delete all entries, that were used to compose the archive entry
		database.Con.Delete(&structs.CheckResultsModel{}, idsToDelete)
	}

	// calulate availabilities etc.
	resultServices := make(map[string]structs.ArchiveResult)

	for name, uds := range udsMap {
		availabilityFull := float64(uds.Ups) / float64(uds.Ups+uds.Downs)
		availability := math.Round(availabilityFull*10000) / 10000 // rounds to four decimal places
		var status structs.Status
		if availability > 0.9 {
			status = structs.Up
		} else if availability > 0.7 {
			status = structs.Problems
		} else {
			status = structs.Down
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
		// DELETE FROM archive_results_models WHERE ID IN (...);
		database.Con.Delete(&older)
	}

}

package schedules

import (
	"database/sql"
	"log"
	"math"
	"time"

	"status-api/database"
	"status-api/structs"
)

// noDowntimeReason is a dummy reason used to mark "no downtime",
// so that downtimes without a reason won't get concatenated if
// there is an error at the first time of checking
const noDowntimeReason = "V2VyIGRhcyBsaWVzdCBpc3QgZG9vZiA7KQ=="

// upsAndDowns counts up-/downtimes of a service and
// collects its downtime intervals (only used locally)
type upsAndDowns struct {
	Ups, Downs int
	Downtimes  []structs.Downtime
}

// StartArchiveTriggerJob starts the goroutine
// to trigger archiving of all checks at 23:59 every day
func StartArchiveTriggerJob() {
	time.Sleep(time.Duration(3) * time.Second)

	ran := make(chan time.Duration)

	job, err := scheduler.Every(1).Day().At("23:59").SingletonMode().Do(func() {
		start := time.Now()
		runArchiving()
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

// aggregate all check results (todays entries; older ones may
// remain if a previous archiving run was missed)
func runArchiving() {
	udsMap := make(map[string]*upsAndDowns)
	lastDownReason := make(map[string]string)

	rows, err := database.Con.Model(&structs.CheckResultsModel{}).Order("id").Rows() // TODO maybe only get todays entries?
	if err != nil {
		panic(err) // should never happen
	}
	idsToDelete := collectResults(rows, udsMap, lastDownReason)

	rows.Close()
	if len(idsToDelete) > 0 { // suppresses error logs
		// delete all entries, that were used to compose the archive entry
		database.Con.Delete(&structs.CheckResultsModel{}, idsToDelete)
	}

	// calulate availabilities etc.
	resultServices := buildArchiveResults(udsMap)

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

// ComputeTodayResults aggregates the check results of the current
// day into an ArchiveResults object. This way, todays downtime
// events can be observed before the archiving job runs at the end
// of the day. Services without any check result today are omitted.
func ComputeTodayResults() *structs.ArchiveResults {
	// aggregate todays check results
	udsMap := make(map[string]*upsAndDowns)
	lastDownReason := make(map[string]string)
	today := dayOnly()

	var todaysChecks []structs.CheckResultsModel
	database.Con.Model(&structs.CheckResultsModel{}).Order("id").Find(&todaysChecks)
	for i := range todaysChecks {
		if todaysChecks[i].Data.At.Before(today) {
			continue // not from today (leftover of a missed archiving run)
		}
		aggregateCheckResults(&todaysChecks[i].Data, udsMap, lastDownReason)
	}

	return &structs.ArchiveResults{
		At:       today,
		Services: buildArchiveResults(udsMap),
	}
}

// collectResults iterates over the given rows of CheckResultsModels,
// aggregates them into udsMap and returns the IDs of all visited rows.
func collectResults(rows *sql.Rows, udsMap map[string]*upsAndDowns, lastDownReason map[string]string) (ids []int) {
	ids = make([]int, 0)
	for rows.Next() { // iterates over CheckResults
		currentCr := &structs.CheckResultsModel{}
		database.Con.ScanRows(rows, currentCr)
		aggregateCheckResults(&currentCr.Data, udsMap, lastDownReason)
		ids = append(ids, int(currentCr.ID))
	}
	return
}

// aggregateCheckResults merges a single CheckResults entry
// (one point in time, all services) into udsMap. Consecutive
// downtimes with the same reason are concatenated.
func aggregateCheckResults(cr *structs.CheckResults, udsMap map[string]*upsAndDowns, lastDownReason map[string]string) {
	for name, service := range cr.Services {

		uds, ok := udsMap[name]
		if !ok { // create the counter, if it does not exist yet
			uds = &upsAndDowns{
				Downtimes: make([]structs.Downtime, 0),
			}
			udsMap[name] = uds
		}
		// initialize an empty lastDownReason with a dummy string
		// (see noDowntimeReason)
		if lastDownReason[name] == "" {
			lastDownReason[name] = noDowntimeReason
		}

		// this is where the aggregation magic happens
		if service.Status == structs.Up {
			uds.Ups++
			// reset lastDownReason to the dummy to prevent the concatenation
			// of downtimes if a previous error occurs at a later time again
			lastDownReason[name] = noDowntimeReason
		} else {
			uds.Downs++
			// check if the previous downtime should be adjusted, or if a new one should be created
			if lastDownReason[name] == service.Reason {
				// if the reason of the last downtime is the same as this one just
				// overwrite the latter 'To' field to the time of the current downtime
				lastDowntime := &uds.Downtimes[len(uds.Downtimes)-1]
				lastDowntime.To = cr.At
			} else {
				// if the reason is a new one, create a new downtime object and overwrite the lastDownReason
				uds.Downtimes = append(uds.Downtimes, structs.Downtime{
					From:   cr.At,
					To:     cr.At,
					Reason: service.Reason,
				})
				lastDownReason[name] = service.Reason
			}
		}
	}
}

// buildArchiveResults calculates the availabilities and
// overall statuses from the aggregated upsAndDowns
func buildArchiveResults(udsMap map[string]*upsAndDowns) map[string]structs.ArchiveResult {
	resultServices := make(map[string]structs.ArchiveResult)
	for name, uds := range udsMap {
		availabilityFull := float64(uds.Ups) / float64(uds.Ups+uds.Downs)
		availability := math.Round(availabilityFull*10000) / 10000 // rounds to four decimal places
		resultServices[name] = structs.ArchiveResult{
			Status:       statusFromAvailability(availability),
			Availability: availability,
			Downtimes:    uds.Downtimes,
		}
	}
	return resultServices
}

// statusFromAvailability maps an availability to an overall
// status (>97% up, >80% problems, otherwise down)
func statusFromAvailability(availability float64) structs.Status {
	if availability > 0.97 {
		return structs.Up
	} else if availability > 0.8 {
		return structs.Problems
	}
	return structs.Down
}

// dayOnly returns the current day at 00:00 local time
func dayOnly() time.Time {
	now := time.Now()
	return time.Date(
		now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location(),
	)
}

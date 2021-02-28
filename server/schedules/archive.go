package schedules

import (
	"fmt"
	"math"
	"time"

	"status-api/database"
	"status-api/structs"
)

func runArchiving(config *structs.Config) {

	dayOnly := func() time.Time {
		now := time.Now()
		return time.Date(
			now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location(),
		)
	}

	// get rows and count up-/downtimes
	rows, err := database.Con.Model(&structs.CheckResultsModel{}).Rows() // TODO maybe only get todays entries?
	if err != nil {
		panic(err)
	}

	type UpsAndDowns struct {
		Ups, Downs int
		Downtimes  []structs.Downtime
	}

	udsMap := make(map[string]*UpsAndDowns) // maps service name to UpsAndDowns counter
	idsToDelete := make([]int, 0)

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
	database.Con.Delete(&structs.CheckResultsModel{}, idsToDelete)

	// Calulate availabilities etc.
	resultServices := make(map[string]structs.ArchiveResult)

	for name, uds := range udsMap {
		availabilityFull := float64(uds.Ups) / float64(uds.Ups+uds.Downs)
		fmt.Println(availabilityFull)
		availability := math.Round(availabilityFull*100) / 100
		fmt.Println(availability)
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

	resultModel := &structs.ArchiveResultsModel{
		Data: structs.ArchiveResults{
			At:       dayOnly(),
			Services: resultServices,
		},
	}

	fmt.Println(resultModel)

	database.Con.Create(resultModel)

}

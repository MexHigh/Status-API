package server

import (
	"net/http"

	"status-api/database"
	"status-api/structs"
)

func latestHandler(w http.ResponseWriter, r *http.Request) {
	st := &structs.CheckResultsModel{}
	database.Con.Last(st)

	respondData(&w, st.Data, 200)
}

func timelineHandler(w http.ResponseWriter, r *http.Request) {
	var tlModel []structs.ArchiveResultsModel
	database.Con.Find(&tlModel)

	tl := make([]structs.ArchiveResults, 0, len(tlModel))
	for _, v := range tlModel {
		tl = append(tl, structs.ArchiveResults{
			At:       v.Data.At,
			Services: v.Data.Services,
		})
	}

	respondData(&w, &tl, 200)
}

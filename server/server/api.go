package server

import (
	"net/http"
	"status-api/database"
	"status-api/structs"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(&w, []byte(`{"response": "pong"}`), 200)
}

func latestHandler(w http.ResponseWriter, r *http.Request) {
	st := &structs.CheckResultsModel{}
	database.Con.Last(st)

	respondInstance(&w, st.Data, 200)
}

func timelineHandler(w http.ResponseWriter, r *http.Request) {

	var tl []structs.ArchiveResultsModel
	database.Con.Find(&tl)

	tlWithoutData := make([]structs.ArchiveResults, 0, len(tl))
	for _, v := range tl {
		tlWithoutData = append(tlWithoutData, structs.ArchiveResults{
			At:       v.Data.At,
			Services: v.Data.Services,
		})
	}

	respondInstance(&w, &tlWithoutData, 200)
}

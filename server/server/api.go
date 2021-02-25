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
	st := &structs.Results{}
	database.Con.Last(st)

	respondInstance(&w, st, 200)
}

func timelineHandler(w http.ResponseWriter, r *http.Request) {

	type Timeline []structs.ArchiveResults
	tl := &Timeline{}
	database.Con.Find(tl)

	respondInstance(&w, tl, 200)
}

package server

import (
	"encoding/json"
	"net/http"
)

func respondInstance(w *http.ResponseWriter, inst interface{}, statusCode int) {
	bytes, err := json.MarshalIndent(inst, "", "    ")
	if err != nil {
		panic(err)
	}
	respondJSON(w, bytes, statusCode)
}

func respondJSON(w *http.ResponseWriter, json []byte, statusCode int) {
	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")
	(*w).WriteHeader(statusCode)
	(*w).Write(json)
}

func respondError(w *http.ResponseWriter, err error) {
	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")
	(*w).WriteHeader(400)
	(*w).Write([]byte(`{"error": "` + err.Error() + `"}`))
}

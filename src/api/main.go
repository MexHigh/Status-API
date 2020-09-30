package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"status-api/checker"
)

const (
	apiHost = "0.0.0.0:3002"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(&w, []byte(`{"status": "ok"}`), 200)
	return
}

func allServicesHandler(w http.ResponseWriter, r *http.Request) {
	allStatus := make(map[string]map[string]string)
	for name, endpoint := range checker.Endpoints {
		status := endpoint.Status
		allStatus[name] = status
	}
	jsonStatus, err := json.MarshalIndent(allStatus, "", "    ")
	if err != nil {
		panic(err)
	}
	respondJSON(&w, jsonStatus, 200)
	return
}

func oneServiceHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["name"]
	if !ok {
		respondError(&w, errors.New("Please provide ?name= parameter"))
		return
	}
	endpoint, err := checker.Endpoints.GetEndpoint(keys[0])
	if err != nil {
		respondError(&w, err)
		return
	}
	jsonStatus, err := endpoint.Status.JSON()
	if err != nil {
		respondError(&w, err)
		return
	}
	respondJSON(&w, jsonStatus, 200)
	return
}

// Start starts the API
func Start() {
	log.Println("Starting status API-server at " + apiHost)
	// endpoints
	router := mux.NewRouter()
	router.HandleFunc("/status", statusHandler).Methods("GET")
	router.HandleFunc("/services", allServicesHandler).Methods("GET")
	router.HandleFunc("/service", oneServiceHandler).Methods("GET")
	// serve
	server := &http.Server{
		Addr:    apiHost,
		Handler: router,
	}
	panic(server.ListenAndServe()) // ListenAndServer never returns a non-nil error
}

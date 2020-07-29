package api

import (
	"log"
	"errors"
	"net/http"
	"github.com/gorilla/mux"
	
	"../checker"
)

const (
    apiHost = "0.0.0.0:3002"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(&w, []byte(`{"status": "ok"}`), 200)
	return
}

func allServicesHandler(w http.ResponseWriter, r *http.Request) {
	jsonStatus, err := checker.Status.JSON()
	if err != nil {
		panic(err)
	}
	respondJSON(&w, jsonStatus, 200)
	return
}

func oneServiceHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["url"]
	if !ok {
		respondError(&w, errors.New("Please provide ?url= parameter"))
		return
	}
	status, err := checker.Status.GetEndpoint(keys[0])
	if err != nil {
		respondError(&w, err)
		return
	}
	jsonStatus, err := status.JSON()
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
			Addr: apiHost,
			Handler: router,
	}
	panic(server.ListenAndServe()) // ListenAndServer never returns a non-nil error
}
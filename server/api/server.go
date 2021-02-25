package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const apiHostDefault = "0.0.0.0:3002"

func pingHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(&w, []byte(`{"status": "ok"}`), 200)
	return
}

// Start starts the API
func Start(apiHost string) error {

	if apiHost == "" {
		log.Println("\"api_host\" not defined in config -> using default")
		apiHost = apiHostDefault
	}

	// endpoints
	router := mux.NewRouter()
	router.HandleFunc("/ping", pingHandler).Methods("GET")

	// serve
	server := &http.Server{
		Addr:    apiHost,
		Handler: router,
	}

	log.Println("API server up at " + apiHost)

	return server.ListenAndServe() // Forever blocking, unless there is an error

}

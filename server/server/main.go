package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const hostDefault = "0.0.0.0:3002"

// Start starts the API and the Frontend server
func Start(host string, serveFrontend bool, frontendPath string) error {

	if host == "" {
		log.Println("\"host\" not defined in config -> using default")
		host = hostDefault
	}

	// endpoints
	router := mux.NewRouter()
	// API router
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/ping", pingHandler).Methods("GET")
	apiRouter.HandleFunc("/services/latest", latestHandler).Methods("GET")
	apiRouter.HandleFunc("/services/timeline", timelineHandler).Methods("GET")
	// frontend router
	if serveFrontend {
		if frontendPath == "" {
			// no "using default" message here
			frontendPath = "../frontend/build"
		}
		router.PathPrefix("/").Handler(
			http.FileServer(http.Dir(frontendPath)),
		)
	}

	// serve
	server := &http.Server{
		Addr:    host,
		Handler: router,
	}

	msg := fmt.Sprintf("Server up at %s serving API", host)
	if serveFrontend {
		msg += " and Frontend"
	}
	log.Println(msg)

	return server.ListenAndServe() // Forever blocking, unless there is an error

}

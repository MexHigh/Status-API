package server

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const hostDefault = "0.0.0.0:3002"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().Unix())
}

// Start starts the API and the Frontend server
func Start(host, frontendPath, dashboardTitle, logoPath string, serveFrontend bool, allowedAPIKeys []string) error {
	if host == "" {
		log.Println("\"host\" not defined in config -> using default")
		host = hostDefault
	}

	if len(allowedAPIKeys) == 0 {
		log.Println("No \"allowed_api_keys\" defined -> genereting a temporary one")

		b := make([]rune, 30)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		newKey := string(b)
		allowedAPIKeys = []string{newKey}

		log.Printf("Your temporary API key is '%s' (keep in mind that it will be regenerated after a restart if you do not generate one by yourself)", newKey)
	}

	// endpoints
	router := mux.NewRouter()
	// API router
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/ping", pingHandler).Methods("GET")
	apiRouter.HandleFunc("/dashboard/title", titleHandlerWith(dashboardTitle)).Methods("GET")
	apiRouter.HandleFunc("/dashboard/logo", imageHandlerWith(logoPath)).Methods("GET")
	// Service status API
	apiRouter.HandleFunc("/services/latest", latestHandler).Methods("GET")
	apiRouter.HandleFunc("/services/timeline", timelineHandler).Methods("GET")
	// message API subrouter (uses authentication)
	messageAPIRouter := apiRouter.NewRoute().Subrouter()
	messageAPIRouter.Use(makeAPIKeyAuthMiddleware(allowedAPIKeys))
	messageAPIRouter.HandleFunc("/messages", rssListMessagesHandler).Methods("GET")
	messageAPIRouter.HandleFunc("/message", rssCreateMessageHandler).Methods("POST")
	messageAPIRouter.HandleFunc("/message/{db_id}", rssChangeMessageHandler).Methods("PATCH")
	messageAPIRouter.HandleFunc("/message/{db_id}", rssDeleteMessageHandler).Methods("DELETE")
	// Atom router
	router.HandleFunc("/messages.atom", rssShowHandler).Methods("GET")
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
		msg += " and frontend"
	}
	log.Println(msg)

	return server.ListenAndServe() // Forever blocking, unless there is an error

}

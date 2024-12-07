package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"
)

// Start starts the API and the Frontend server
func Start(host, frontendPath, dashboardTitle, logoPath string, serveFrontend bool, allowedAPIKeys []string) error {
	// create store for rate limiting
	store, err := memorystore.New(&memorystore.Config{
		Tokens:   10,
		Interval: 10 * time.Second,
	})
	if err != nil {
		return err
	}
	limiter, err := httplimit.NewMiddleware(store, httplimit.IPKeyFunc())
	if err != nil {
		return err
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
	// Auth API
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	authRouter.Use(limiter.Handle)
	authRouter.HandleFunc("/api-key/test", makeAPIKeyAuthOkHandler(allowedAPIKeys)).Methods("POST")
	// message API subrouter (uses authentication)
	messageAPIRouter := apiRouter.NewRoute().Subrouter()
	messageAPIRouter.Use(makeAPIKeyAuthMiddleware(allowedAPIKeys))
	messageAPIRouter.Use(limiter.Handle)
	messageAPIRouter.HandleFunc("/messages", rssListMessagesHandler).Methods("GET")
	messageAPIRouter.HandleFunc("/message", rssCreateMessageHandler).Methods("POST")
	messageAPIRouter.HandleFunc("/message/{db_id}", rssChangeMessageHandler).Methods("PATCH")
	messageAPIRouter.HandleFunc("/message/{db_id}", rssDeleteMessageHandler).Methods("DELETE")
	// Atom router
	router.HandleFunc("/messages.atom", rssShowHandler).Methods("GET")
	// frontend router
	if serveFrontend {
		addReactRoute := func(path string) {
			router.PathPrefix(path).Handler(http.StripPrefix(path, http.FileServer(http.Dir(frontendPath))))
		}
		addReactRoute("/admin")
		addReactRoute("/")
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

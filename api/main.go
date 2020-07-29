package api

import (
        "fmt"
        "net/http"
        "github.com/gorilla/mux"
)

const (
        apiHost = "0.0.0.0:3002"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
        respondJSON(&w, []byte(`{"status": "ok"}`), 200)
        return
}

// Start starts the API
func Start() {
        fmt.Println("Starting Command Quest Server at " + apiHost)
        // endpoints
        router := mux.NewRouter()
        router.HandleFunc("/status", statusHandler).Methods("GET")
        // serve
        server := &http.Server{
                Addr: apiHost,
                Handler: router,
        }
        fmt.Println("Up")
        panic(server.ListenAndServe()) // ListenAndServer never returns a non-nil error
}
package server

import (
	"net/http"
	"os"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	respondData(&w, "pong", 200)
}

func titleHandlerWith(title string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if title == "" {
			respondData(&w, "Service Status", http.StatusOK)
		} else {
			respondData(&w, title, http.StatusOK)
		}
	}
}

func imageHandlerWith(imagePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open(imagePath)
		if err != nil {
			respondError(&w, err, http.StatusInternalServerError)
			return
		}
		defer file.Close()
		stat, _ := file.Stat()

		http.ServeContent(w, r, stat.Name(), stat.ModTime(), file)
	}
}

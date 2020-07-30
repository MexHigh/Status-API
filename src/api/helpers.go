package api

import (
	"net/http"
)

// checks if a method ("allowedMethods") was used
func method(r *http.Request, allowedMethods ...string) bool {
	allowed := true
	for _, method := range allowedMethods {
		if r.Method != method {
			allowed = false
		}
	}
	return allowed
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

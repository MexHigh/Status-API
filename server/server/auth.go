package server

import (
	"errors"
	"net/http"
)

// makeAPIKeyAuthOkHandler creates a http.HandlerFunc which checks, if the
// provided API key in the x-api-key header is valid.
func makeAPIKeyAuthOkHandler(allowedAPIKeys []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("X-API-Key")
		if authToken == "" {
			respondError(&w, errors.New("header 'X-API-Key' missing"), 400)
			return
		}
		credsMatch := false
		for _, token := range allowedAPIKeys {
			if authToken == token {
				credsMatch = true
				break
			}
		}
		if credsMatch {
			respondData(&w, "ok", 200)
		} else {
			respondError(&w, errors.New("forbidden"), 403)
		}
	}
}

func makeAPIKeyAuthMiddleware(allowedAPIKeys []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqAuthToken := r.Header.Get("X-API-Key")

			// check if token is provided
			if reqAuthToken == "" {
				respondError(&w, errors.New("unauthorized"), 401)
				return
			}

			// check if provided token is valid
			credsMatch := false
			for _, token := range allowedAPIKeys {
				if reqAuthToken == token {
					credsMatch = true
					break
				}
			}
			if !credsMatch {
				respondError(&w, errors.New("forbidden"), 403)
				return
			}

			// continue with middleware chain
			next.ServeHTTP(w, r)
		})
	}
}

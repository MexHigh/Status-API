package server

import (
	"encoding/json"
	"net/http"
)

func respondPlainJSON(w *http.ResponseWriter, json []byte, statusCode int) {
	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")
	(*w).WriteHeader(statusCode)
	(*w).Write(json)
}

type Response struct {
	Error    *string     `json:"error"`
	Response interface{} `json:"response,omitempty"`
}

func respondData(w *http.ResponseWriter, inst interface{}, statusCode int) {
	wrapped := Response{nil, inst}
	bytes, err := json.MarshalIndent(wrapped, "", "    ")
	if err != nil {
		panic(err)
	}
	respondPlainJSON(w, bytes, statusCode)
}

func respondError(w *http.ResponseWriter, err error, statusCode int) {
	errString := err.Error()
	wrapped := Response{&errString, nil}
	bytes, err := json.MarshalIndent(wrapped, "", "    ")
	if err != nil {
		panic(err)
	}
	respondPlainJSON(w, bytes, statusCode)
}

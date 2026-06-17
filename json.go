package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeErrorResponse(w http.ResponseWriter, code int, msg string) {
	type responseBody struct {
		Error string `json:"error"`
	}

	response := responseBody{
		Error: msg,
	}

	message, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(message)
}

func writeJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

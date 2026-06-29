package main

import (
	"log"
	"net/http"
)

func (apiCfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {

	if apiCfg.platform != "dev" {
		w.Header().Add("Content-type", "text/plain; charset=utf-8")
		w.WriteHeader(403)
		w.Write([]byte("FORBIDDEN"))
		return
	}

	if err := apiCfg.dbQueries.DeleteAllChirps(r.Context()); err != nil {
		log.Printf("failed to reset chirps: %v", err)
		w.Header().Add("Content-type", "text/plain; charset=utf-8")
		w.WriteHeader(500)
		w.Write([]byte("FAILED TO RESET"))
		return
	}

	if err := apiCfg.dbQueries.DeleteAllUsers(r.Context()); err != nil {
		log.Printf("failed to reset users: %v", err)
		w.Header().Add("Content-type", "text/plain; charset=utf-8")
		w.WriteHeader(500)
		w.Write([]byte("FAILED TO RESET"))
		return
	}

	w.Header().Add("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("RESET"))
}

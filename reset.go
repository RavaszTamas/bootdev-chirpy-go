package main

import "net/http"

func (apiCfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {

	if apiCfg.platform != "dev" {
		w.Header().Add("Content-type", "text/plain; charset=utf-8")
		w.WriteHeader(403)
		w.Write([]byte("FORBIDDEN"))
	}

	apiCfg.dbQueries.DeleteAllUsers(r.Context())
	apiCfg.dbQueries.DeleteAllChirps(r.Context())

	w.Header().Add("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("RESET"))
}

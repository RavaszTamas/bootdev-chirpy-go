package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (apiCfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	type requestData struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)

	data := requestData{}

	err := decoder.Decode(&data)

	if err != nil {
		log.Printf("Error decoding request body %s", err)
		writeErrorResponse(w, 500, "Something went wrong")
		return
	}

	log.Printf("Obtained message: %s", data.Email)

	user, err := apiCfg.dbQueries.CreateUser(r.Context(), data.Email)

	if err != nil {
		log.Printf("Failed to create user: %v", err)
		writeErrorResponse(w, 500, "Failed to create user")
		return
	}

	writeJSONResponse(w, 201, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}

package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/RavaszTamas/bootdev-chirpy-go/internal/auth"
	"github.com/RavaszTamas/bootdev-chirpy-go/internal/database"
)

func (apiCfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	type requestData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)

	data := requestData{}

	err := decoder.Decode(&data)

	if err != nil {
		log.Printf("Error decoding request body %s", err)
		writeErrorResponse(w, 500, "Something went wrong")
		return
	}

	hashedPassword, err := auth.HashPassword(data.Password)

	if err != nil {
		log.Printf("Failed to create user: %v", err)
		writeErrorResponse(w, 500, "Failed to create user")
		return
	}

	user, err := apiCfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email: data.Email,
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
	})

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

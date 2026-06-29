package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RavaszTamas/bootdev-chirpy-go/internal/auth"
)

type message struct {
	Message string `json:"message"`
}

func (apiCfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
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

	user, err := apiCfg.dbQueries.GetUserByEmail(r.Context(), data.Email)

	if err != nil {
		log.Printf("Failed to login: %v", err)
		writeErrorResponse(w, 500, "Failed to login")
		return
	}

	value, err := auth.CheckPasswordHash(data.Password, user.HashedPassword.String)

	if err != nil {
		log.Printf("Failed to login: %v", err)
		writeErrorResponse(w, 500, "Failed to login")
		return
	}

	if !value {
		writeJSONResponse(w, 401, message{
			Message: "Invalid password!",
		})
		return
	}

	writeJSONResponse(w, 200, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/RavaszTamas/bootdev-chirpy-go/internal/auth"
)

type message struct {
	Message string `json:"message"`
}

const defaultExpiration = 3600

func (apiCfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type requestData struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)

	data := requestData{}

	err := decoder.Decode(&data)

	if err != nil {
		log.Printf("Error decoding request body %s", err)
		writeErrorResponse(w, 500, "Something went wrong")
		return
	}

	expiresIn := defaultExpiration

	if data.ExpiresInSeconds > 0 && data.ExpiresInSeconds < defaultExpiration {
		expiresIn = data.ExpiresInSeconds
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

	token, err := auth.MakeJWT(user.ID, apiCfg.tokenSecret, time.Duration(expiresIn)*time.Second)

	if err != nil {
		log.Printf("Failed to login: %v", err)
		writeErrorResponse(w, 500, "Failed to login")
		return
	}

	writeJSONResponse(w, 200, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})

}

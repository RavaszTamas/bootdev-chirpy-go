package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/RavaszTamas/bootdev-chirpy-go/internal/auth"
	"github.com/RavaszTamas/bootdev-chirpy-go/internal/database"
)

type message struct {
	Message string `json:"message"`
}

const defaultExpiration = 3600
const refreshTokenExpiration = 60 * 24 * 60 * 60 * time.Second

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

	token, err := auth.MakeJWT(user.ID, apiCfg.tokenSecret, time.Duration(defaultExpiration)*time.Second)

	if err != nil {
		log.Printf("Failed to login: %v", err)
		writeErrorResponse(w, 500, "Failed to login")
		return
	}

	refresh_token, err := apiCfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     auth.MakeRefreshToken(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Duration(refreshTokenExpiration)),
	})

	if err != nil {
		log.Printf("Failed to create refresh token: %v", err)
		writeErrorResponse(w, 500, "Failed to login")
		return
	}

	writeJSONResponse(w, 200, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refresh_token.Token,
	})

}

func (apiCfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		log.Printf("Failed to get refresh token: %v", err)
		writeErrorResponse(w, 401, "Missing refresh token")
		return
	}

	user, err := apiCfg.dbQueries.GetUserFromRefreshToken(r.Context(), token)

	if err != nil {
		log.Printf("Failed to get user: %v", err)
		writeErrorResponse(w, 401, "Invalid refresh token")
		return
	}

	jwtToken, err := auth.MakeJWT(user.ID, apiCfg.tokenSecret, time.Duration(defaultExpiration)*time.Second)

	if err != nil {
		log.Printf("Failed to login: %v", err)
		writeErrorResponse(w, 500, "Failed to login")
		return
	}

	writeJSONResponse(w, 200, Token{
		Token: jwtToken,
	})

}

func (apiCfg *apiConfig) handleRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		log.Printf("Failed to get refresh token: %v", err)
		writeErrorResponse(w, 401, "Missing refresh token")
		return
	}

	_, err = apiCfg.dbQueries.RevokeRefreshToken(r.Context(), token)

	if err != nil {
		log.Printf("Failed to get user: %v", err)
		writeErrorResponse(w, 500, "Failed to revoke refresh token")
		return
	}

	writeJSONResponse(w, 204, nil)

}

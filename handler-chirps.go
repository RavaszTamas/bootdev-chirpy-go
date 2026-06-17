package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/RavaszTamas/bootdev-chirpy-go/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerChirpGetOne(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")

	id, err := uuid.Parse(chirpID)

	if err != nil {
		log.Printf("Failed to get parameter chirpID: %s. %v", chirpID, err)
		writeErrorResponse(w, 400, fmt.Sprintf("Invalid parameter for chirpID: %s. Expecting UUID", chirpID))
		return
	}

	chirp, err := apiCfg.dbQueries.GetChirpById(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Chirp not found: %s", id)
			writeErrorResponse(w, 404, fmt.Sprintf("Chirp not found: %s", id))
			return
		}

		log.Printf("Database error: %v", err)
		writeErrorResponse(w, 500, "internal server error")
		return
	}

	writeJSONResponse(w, 200, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func (apiCfg *apiConfig) handlerChirpsGetAll(w http.ResponseWriter, r *http.Request) {
	chirps, err := apiCfg.dbQueries.GetAllChirps(r.Context())

	if err != nil {
		log.Printf("Failed to create user: %v", err)
		writeErrorResponse(w, 500, "Failed to get all chirps")
		return
	}

	chirps_response := make([]Chirp, 0, len(chirps))

	for _, elem := range chirps {
		chirps_response = append(chirps_response, Chirp{
			ID:        elem.ID,
			CreatedAt: elem.CreatedAt,
			UpdatedAt: elem.UpdatedAt,
			Body:      elem.Body,
			UserId:    elem.UserID,
		})
	}

	writeJSONResponse(w, 200, chirps_response)
}

func (apiCfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type requestData struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)

	data := requestData{}

	err := decoder.Decode(&data)

	if err != nil {
		log.Printf("Error decoding request body %s", err)
		writeErrorResponse(w, 500, "Something went wrong")
		return
	}

	log.Printf("Obtained message %s", data.Body)

	if len(data.Body) > 140 {
		log.Printf("Chirp is too long %d", len(data.Body))
		writeErrorResponse(w, 400, "Chirp is too long")
		return
	}

	chirp, err := apiCfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   data.Body,
		UserID: data.UserId,
	})

	if err != nil {
		log.Printf("Error decoding request body %s", err)
		writeErrorResponse(w, 500, "Something went wrong")
		return
	}

	writeJSONResponse(w, 201, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		UserId:    chirp.UserID,
		Body:      chirp.Body,
	})

}

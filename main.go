package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/RavaszTamas/bootdev-chirpy-go/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func main() {
	// -- env
	godotenv.Load()

	platform := os.Getenv("PLATFORM")

	// -- db setup

	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatalf("Failed to open server connection: %v", err)
	}

	dbQueries := database.New(db)

	// -- server setuo

	mux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      dbQueries,
		platform:       platform,
	}
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)

	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsGetAll)

	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpGetOne)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUserCreate)

	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)

	server := http.Server{
		Handler: mux,
		Addr:    "localhost:8080",
	}

	server.ListenAndServe()
}

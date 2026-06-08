package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/RavaszTamas/bootdev-chirpy-go/internal/database"
	"github.com/RavaszTamas/bootdev-chirpy-go/validation"
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

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

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

func writeValidResponse(body string, w http.ResponseWriter) {
	type responseBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	body = validation.ReplaceBadWords(body)

	response := responseBody{
		CleanedBody: body,
	}

	message, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(message)
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

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/html")
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
		`, apiCfg.fileserverHits.Load())))
	})

	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {

		if apiCfg.platform != "dev" {
			w.Header().Add("Content-type", "text/plain; charset=utf-8")
			w.WriteHeader(403)
			w.Write([]byte("FORBIDDEN"))
		}

		apiCfg.dbQueries.DeleteAllUsers(r.Context())

		w.Header().Add("Content-type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("RESET"))
	})

	mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		type requestData struct {
			Body string `json:"body"`
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
		} else {
			writeValidResponse(data.Body, w)
		}

	})

	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
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

		responseUser := User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}

		message, err := json.Marshal(responseUser)

		if err != nil {
			log.Printf("Failed to marshal user response %v", err)
			writeErrorResponse(w, 500, "Failed to marshal response")
			return
		}

		w.WriteHeader(201)
		w.Header().Set("Content-Type", "application/json")
		w.Write(message)

	})

	server := http.Server{
		Handler: mux,
		Addr:    "localhost:8080",
	}

	server.ListenAndServe()
}

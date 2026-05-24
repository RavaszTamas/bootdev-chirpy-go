package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/RavaszTamas/bootdev-chirpy-go/validation"
)

type apiConfig struct {
	fileserverHits atomic.Int32
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
	mux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
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
		}

		log.Printf("Obtained message %s", data.Body)

		if len(data.Body) > 140 {
			log.Printf("Chirp is too long %d", len(data.Body))
			writeErrorResponse(w, 400, "Chirp is too long")
		} else {
			writeValidResponse(data.Body, w)
		}

	})

	server := http.Server{
		Handler: mux,
		Addr:    "localhost:8080",
	}

	server.ListenAndServe()
}

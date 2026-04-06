package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/NinjaCrusader/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	cfgPlatform    string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("Platform is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("there was an issue opening the db connection: %v\n", err)
		return
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		db:          dbQueries,
		cfgPlatform: platform,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc((http.StripPrefix("/app", http.FileServer(http.Dir("."))))))
	mux.HandleFunc("POST /api/login", apiCfg.handlerUserLogin)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetSingleChirp)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handlerValidate)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	s := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	listenErr := s.ListenAndServe()
	if listenErr != nil {
		log.Fatalf("there was an error with the server: %v\n", err)
	}
}

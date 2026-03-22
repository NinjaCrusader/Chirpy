package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {

	count := cfg.fileserverHits.Load()

	adminMetrics := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>\n", count)

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	_, err := w.Write([]byte(adminMetrics))
	if err != nil {
		log.Printf("there was an error writing the hit count: %v\n", err)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {

	count := fmt.Sprintf("Hits: %v\n", cfg.fileserverHits.Load())

	_, err := w.Write([]byte(count))
	if err != nil {
		log.Printf("there was an error writing the hit count: %v\n", err)
	}
}

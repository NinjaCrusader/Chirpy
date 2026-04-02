package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/NinjaCrusader/Chirpy/internal/database"
	"github.com/google/uuid"
)

type chirpPayload struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type result struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {

	chirp := chirpPayload{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("there was an issue decoding the chirp request: %v\n", err)
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, http.StatusInternalServerError, "Chirp is too long")
		return
	}

	cleaned := removeBadWords(chirp.Body)

	insertChripParam := database.InsertChirpParams{
		Body:   cleaned,
		UserID: chirp.UserID,
	}

	insertChirp, err := cfg.db.InsertChirp(r.Context(), insertChripParam)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("there was an issue with inserting chirp into db: %v\n", err)
		return
	}

	createdResult := result{
		ID:        insertChirp.ID,
		CreatedAt: insertChirp.CreatedAt,
		UpdatedAt: insertChirp.UpdateAt,
		Body:      insertChirp.Body,
		UserID:    insertChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, createdResult)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	final := []result{}

	chirps, err := cfg.db.GetUsers(r.Context())
	if err != nil {
		log.Printf("there was an error getting chirps from the db: %v\n")
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	for i := 0; i < len(chirps); i++ {
		row := chirps[i]
		chirp := result{
			ID:        row.ID,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdateAt,
			Body:      row.Body,
			UserID:    row.UserID,
		}

		final = append(final, chirp)
	}

	respondWithJSON(w, http.StatusOK, final)

}

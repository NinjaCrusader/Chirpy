package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type requestParams struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := requestParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("There was an error decoding the request: %v\n", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		if dbErr, ok := err.(*pq.Error); ok {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			log.Printf("There was an error with postgres creating the user: %v\n", dbErr.Code)
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
			log.Printf("There was an error creating the user: %v\n", err)
			return
		}
	}

	userParams := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusCreated, userParams)

}

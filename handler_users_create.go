package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NinjaCrusader/Chirpy/internal/auth"
	"github.com/NinjaCrusader/Chirpy/internal/database"
	"github.com/lib/pq"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := requestParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("There was an error decoding the request: %v\n", err)
		return
	}

	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("something went wrong while hashing the password: %v\n", err)
		return
	}

	createUserParam := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed,
	}

	user, err := cfg.db.CreateUser(r.Context(), createUserParam)
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

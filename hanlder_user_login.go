package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NinjaCrusader/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {

	decode := json.NewDecoder(r.Body)
	params := requestParams{}
	err := decode.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("there was an error getting the request information: %v\n", err)
		return
	}

	user, err := cfg.db.FindUserToLogin(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		log.Printf("there was an error finding the user to login: %v\n", err)
		return
	}

	validCheck, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || validCheck == false {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		log.Printf("there was an error finding the user to login: %v\n", err)
		return
	}

	res := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusOK, res)
}

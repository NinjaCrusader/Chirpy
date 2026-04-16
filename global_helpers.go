package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type requestParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	Expires  int    `json:"expires_in_seconds"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {

	type returnVals struct {
		Error string `json:"error"`
	}

	respBody := returnVals{
		Error: msg,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %v\n", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling response JSON: %v\n", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func removeBadWords(s string) string {

	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}

	split := strings.Split(s, " ")

	for i := 0; i < len(split); i++ {
		copy := split[i]
		if slices.Contains(bannedWords, strings.ToLower(copy)) {
			split[i] = "****"
			continue
		}
	}

	final := strings.Join(split, " ")

	return final
}

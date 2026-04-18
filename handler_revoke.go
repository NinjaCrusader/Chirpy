package main

import (
	"log"
	"net/http"

	"github.com/NinjaCrusader/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {

	headerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("there was an error getting the header token when revoking: %v\n", err)
		return
	}

	_, dbErr := cfg.db.RevokeRefreshToken(r.Context(), headerToken)
	if dbErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("there was an error revoking the refresh Token: %v\n", dbErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

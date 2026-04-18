package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/NinjaCrusader/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {

	type resToken struct {
		Token string `json:"token"`
	}

	headerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("something went wrong getting the headers for refreshing the refresh token: %v\n", err)
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(context.Background(), headerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Something went wrong")
		log.Printf("something went wrong getting the refresh token from the db: %v\n", err)
		return
	} else if time.Now().After(refreshToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Something went wrong")
		log.Printf("the token is already expired")
		return
	} else if refreshToken.RevokedAt.Valid == true {
		respondWithError(w, http.StatusUnauthorized, "Something went wrong")
		log.Printf("token was revoked")
		return
	}

	newToken, err := auth.MakeJWT(refreshToken.UserID, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("there was an error creating the refreshed jwt: %v\n", err)
		return
	}

	res := resToken{
		Token: newToken,
	}

	respondWithJSON(w, http.StatusOK, res)

}

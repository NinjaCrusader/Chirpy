package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {

	if cfg.cfgPlatform != "dev" {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	type resetSuccessParams struct {
		Body string `json:"body"`
	}

	success := resetSuccessParams{
		Body: "Success",
	}

	removeUsersErr := cfg.db.RemoveUsers(r.Context())
	if removeUsersErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("there was an error removeing users from the users table: %v\n", removeUsersErr)
		return
	}
	respondWithJSON(w, http.StatusOK, success)

}

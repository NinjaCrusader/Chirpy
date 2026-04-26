package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type polkaData struct {
	UserID uuid.UUID `json:"user_id"`
}

type polkaWebhookRequest struct {
	Event string    `json:"event"`
	Data  polkaData `json:"data"`
}

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {

	webhook := polkaWebhookRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&webhook)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		log.Printf("there was an error while trying to decode the polka webhook request: %v\n", err)
		return
	}

	if webhook.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if _, err := cfg.db.UpgradeUser(r.Context(), webhook.Data.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Something went wrong")
			log.Printf("no user could be found to upgrade: %v\n", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		log.Printf("something went wrong when trying to upgrade user: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

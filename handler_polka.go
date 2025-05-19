package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jdjaxon/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolka(w http.ResponseWriter, r *http.Request) {
	type polkaWebhookReq struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	var req polkaWebhookReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			err,
			"Error decoding webhook parameters",
		)
		return
	}

	requestingApiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, "Could not retrieve token")
		return
	}

	if requestingApiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, err, "Unauthorized request")
		return
	}

	if req.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	err = cfg.db.UpgradeUserToChirpyRed(r.Context(), req.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, err, "User not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, err, "Could not upgrade user")
		}

		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

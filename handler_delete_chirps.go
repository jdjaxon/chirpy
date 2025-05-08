package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jdjaxon/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirpByID(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, "Could not retrieve token")
		return
	}

	requestingUserID, err := auth.ValidateJWT(tokenString, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, "Invalid token")
		return
	}

	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			err,
			"Failed to parse UUID",
		)
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), uuid.UUID(chirpUUID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(
				w,
				http.StatusNotFound,
				err,
				"Chirp not found",
			)
		} else {
			respondWithError(
				w,
				http.StatusInternalServerError,
				err,
				"Failed to retrieve chirp",
			)
		}

		return
	}

	if chirp.UserID != requestingUserID {
		respondWithError(w, http.StatusForbidden, err, "Action not allowed")
		return
	}

	err = cfg.db.DeleteChirpByID(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err, "Could not delete chirp")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

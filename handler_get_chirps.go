package main

import (
	"database/sql"
	"errors"
	"net/http"
	"slices"

	"github.com/google/uuid"
	"github.com/jdjaxon/chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorIDStr := r.URL.Query().Get("author_id")
	sortStr := r.URL.Query().Get("sort")

	var chirps []database.Chirp
	var err error

	if authorIDStr != "" {
		authorUUID, err := uuid.Parse(authorIDStr)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Failed to parse author UUID")
			return
		}
		chirps, err = cfg.db.GetChirpsByAuthor(r.Context(), authorUUID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Failed to get chirps")
			return
		}
	} else {
		chirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err, "Failed to get chirps")
			return
		}
	}

	// Chirps are sorted by created_at in ascending order from DB.
	if sortStr == "desc" {
		slices.Reverse(chirps)
	}

	type response struct {
		Chirps []Chirp
	}

	resp := []Chirp{}
	for _, chirp := range chirps {
		resp = append(resp, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) handlerGetChirpsByID(w http.ResponseWriter, r *http.Request) {
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

	resp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

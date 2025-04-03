package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			err,
			"Failed to get chirps",
		)
		return
	}

	type response struct {
		Chirps []Chirp
	}
	resp := response{
		Chirps: []Chirp{},
	}
	for _, chirp := range chirps {
		resp.Chirps = append(resp.Chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusCreated, resp)
}

package main

import (
	"net/http"

	"github.com/jdjaxon/chirpy/internal/auth"
)

// handlerRevoke -
func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, "Could not retrieve token")
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, "Invalid token")
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

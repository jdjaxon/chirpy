package main

import (
	"net/http"
	"time"

	"github.com/jdjaxon/chirpy/internal/auth"
	"github.com/jdjaxon/chirpy/internal/database"
)

// handlerRefresh -
func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, "Could not retrieve token")
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, "Invalid token")
		return
	}

	if !isRefreshTokenValid(refreshToken) {
		respondWithError(w, http.StatusUnauthorized, err, "Invalid token")
		return
	}

	newJWT, err := auth.MakeJWT(refreshToken.UserID, cfg.tokenSecret, cfg.jwtTTL)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, "Invalid token")
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	resp := response{
		Token: newJWT,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

// isRefreshTokenValid -
func isRefreshTokenValid(t database.RefreshToken) bool {
	return !t.RevokedAt.Valid && time.Now().Before(t.ExpiresAt)
}

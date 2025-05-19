package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jdjaxon/chirpy/internal/auth"
	"github.com/jdjaxon/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			err,
			"Error decoding request parameters",
		)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			err,
			"Incorrect email or password",
		)
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, req.Password)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			err,
			"Incorrect email or password",
		)
		return
	}

	userToken, err := auth.MakeJWT(user.ID, cfg.tokenSecret, cfg.jwtTTL)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			err,
			"Failed to create token",
		)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			err,
			"Failed to create refresh token",
		)
		return
	}

	err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(cfg.refreshTokenTTL),
	})
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			err,
			"Failed to created refresh token",
		)
		return
	}

	resp := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        userToken,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

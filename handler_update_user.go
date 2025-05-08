package main

import (
	"encoding/json"
	"net/http"

	"github.com/jdjaxon/chirpy/internal/auth"
	"github.com/jdjaxon/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, "Could not retrieve token")
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, "Invalid token")
		return
	}

	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	err = decoder.Decode(&req)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			err,
			"Error decoding request parameters",
		)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			err,
			"Error hashing password",
		)
		return
	}

	refreshToken, err := cfg.db.GetRefreshTokenByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err, "Invalid token")
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	})
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			err,
			"Error updating user email or password",
		)
		return
	}

	resp := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        tokenString,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

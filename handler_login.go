package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jdjaxon/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds,omitempty"`
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

	expiresIn := time.Second * 3600
	if req.ExpiresInSeconds != nil {
		reqExpiry := time.Duration(*req.ExpiresInSeconds) * time.Second
		if reqExpiry <= time.Second*3600 {
			expiresIn = reqExpiry
		}
	}

	// MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error)
	userToken, err := auth.MakeJWT(user.ID, cfg.tokenSecret, expiresIn)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			err,
			"Failed to create token",
		)
		return
	}

	resp := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     userToken,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

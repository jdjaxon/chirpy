package main

import (
	"encoding/json"
	"net/http"

	"github.com/jdjaxon/chirpy/internal/auth"
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

	resp := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

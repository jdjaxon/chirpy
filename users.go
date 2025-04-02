package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type newUserReq struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	userReq := newUserReq{}
	err := decoder.Decode(&userReq)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			err,
			"Error decoding request parameters",
		)
		return
	}

	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	user, err := cfg.db.CreateUser(r.Context(), userReq.Email)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			err,
			"Error creating user",
		)
		return
	}

	resp := User(user)
	respondWithJSON(w, http.StatusCreated, resp)
}

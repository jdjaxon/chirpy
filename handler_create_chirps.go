package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jdjaxon/chirpy/internal/auth"
	"github.com/jdjaxon/chirpy/internal/database"
)

const maxChirpLen = 140

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Body string `json:"body"`
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

	if len(req.Body) > maxChirpLen {
		respondWithError(w, http.StatusBadRequest, nil, "Chirp is too long")
		return
	}
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleanedBody := cleanBody(req.Body, badWords)

	newChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	})
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			err,
			"Error creating chirp",
		)
		return
	}

	resp := Chirp{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	}
	respondWithJSON(w, http.StatusCreated, resp)
}

func cleanBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		if _, exists := badWords[strings.ToLower(word)]; exists {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

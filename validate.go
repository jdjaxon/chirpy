package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

const maxChirpLen = 140

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	chirp := request{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			err,
			"Error decoding request parameters",
		)
		return
	}

	if len(chirp.Body) > maxChirpLen {
		respondWithError(w, http.StatusBadRequest, nil, "Chirp is too long")
		return
	}

	type response struct {
		Body string `json:"cleaned_body"`
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleanedBody := cleanBody(chirp.Body, badWords)
	resp := response{
		Body: cleanedBody,
	}
	respondWithJSON(w, http.StatusOK, resp)
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

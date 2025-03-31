package main

import (
	"encoding/json"
	"net/http"
)

const maxChirpLen = 140

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirpReq struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	chirp := chirpReq{}
	err := decoder.Decode(&chirp)
	if err != nil {
		writeJSONWithError(
			w,
			http.StatusInternalServerError,
			err,
			"Error decoding request parameters",
		)
		return
	}

	if len(chirp.Body) > maxChirpLen {
		writeJSONWithError(w, http.StatusBadRequest, nil, "Chirp is too long")
		return
	}

	type validResp struct {
		Valid bool `json:"valid"`
	}
	resp := validResp{
		Valid: true,
	}
	writeJSON(w, http.StatusOK, resp)
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, statusCode int, err error, msg string) {
	if err != nil {
		log.Printf("%v: %v\n", err.Error(), msg)
	}
	type errorResp struct {
		Error string `json:"error"`
	}
	resp := errorResp{
		Error: msg,
	}
	respondWithJSON(w, statusCode, resp)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("content-type", "application/json")
	jsonData, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(jsonData)
}

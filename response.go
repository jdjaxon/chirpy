package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJSONWithError(w http.ResponseWriter, statusCode int, err error, msg string) {
	if err != nil {
		log.Println(err)
	}
	type errorResp struct {
		Error string `json:"error"`
	}
	resp := errorResp{
		Error: msg,
	}
	writeJSON(w, statusCode, resp)
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("content-type", "application/json")
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(jsonData)
}

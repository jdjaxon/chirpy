package main

import (
	"log"
	"net/http"
)

func main() {
	serverRootDir := "."
	port := "8080"
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(serverRootDir))))
	mux.HandleFunc("/healthz", handlerHealthcheck)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", serverRootDir, port)
	server.ListenAndServe()
}

func handlerHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusOKText))
}

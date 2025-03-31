package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	serverRootDir := "."
	port := "8080"
	mux := http.NewServeMux()
	cfg := &apiConfig{}
	fileserverHanlder := http.StripPrefix(
		"/app",
		http.FileServer(http.Dir(serverRootDir)),
	)

	mux.Handle("/app/", cfg.middlewareMetricsInc(fileserverHanlder))
	mux.HandleFunc("GET /api/healthz", handlerHealthcheck)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerResetMetrics)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", serverRootDir, port)
	log.Fatal(server.ListenAndServe())
}

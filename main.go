package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

const appEndpoint = "/app"

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	serverRootDir := "."
	port := "8080"
	mux := http.NewServeMux()
	cfg := &apiConfig{}

	mux.Handle(
		appEndpoint+"/",
		cfg.middlewareMetricsInc(
			http.StripPrefix(appEndpoint, http.FileServer(http.Dir(serverRootDir))),
		),
	)
	mux.HandleFunc("GET /api/healthz", handlerHealthcheck)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerResetMetrics)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", serverRootDir, port)
	log.Fatal(server.ListenAndServe())
}

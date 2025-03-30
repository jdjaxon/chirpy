package main

import (
	"fmt"
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
	mux.HandleFunc("/healthz", handlerHealthcheck)
	mux.HandleFunc("/metrics", cfg.handlerGetMetrics)
	mux.HandleFunc("/reset", cfg.handlerResetMetrics)

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
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerGetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

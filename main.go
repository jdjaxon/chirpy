package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/jdjaxon/chirpy/internal/database"
	"github.com/joho/godotenv"
	// You have to import the driver, but it don't use it directly.
	// The underscore tells Go that you're importing it for its side effects.
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is required")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM is required")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Failed to connect to DB: %s\n", err)
		return
	}
	dbQueries := database.New(dbConn)

	serverRootDir := "."
	port := "8080"

	cfg := &apiConfig{
		db:       dbQueries,
		platform: platform,
	}

	fileserverHanlder := http.StripPrefix(
		"/app",
		http.FileServer(http.Dir(serverRootDir)),
	)

	mux := http.NewServeMux()

	mux.Handle("/app/", cfg.middlewareMetricsInc(fileserverHanlder))

	mux.HandleFunc("GET /api/healthz", handlerHealthcheck)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("POST /api/users", cfg.handlerUsers)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", serverRootDir, port)
	log.Fatal(server.ListenAndServe())
}

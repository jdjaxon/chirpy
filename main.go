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
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Failed to connect to DB: %s\n", err)
		return
	}
	dbQueries := database.New(db)

	serverRootDir := "."
	port := "8080"
	mux := http.NewServeMux()
	cfg := &apiConfig{
		db: dbQueries,
	}

	fileserverHanlder := http.StripPrefix(
		"/app",
		http.FileServer(http.Dir(serverRootDir)),
	)

	mux.Handle("/app/", cfg.middlewareMetricsInc(fileserverHanlder))

	mux.HandleFunc("GET /api/healthz", handlerHealthcheck)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("POST /api/users", cfg.handlerUsers)

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerResetMetrics)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", serverRootDir, port)
	log.Fatal(server.ListenAndServe())
}

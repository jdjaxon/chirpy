package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/jdjaxon/chirpy/internal/database"
	"github.com/joho/godotenv"
	// You have to import the driver, but it don't use it directly.
	// The underscore tells Go that you're importing it for its side effects.
	_ "github.com/lib/pq"
)

var ErrFailedToLoadEnvVars = errors.New("failed to load environment variables")
var ErrDatabaseConnFailed = errors.New("failed to connect to database")

type apiConfig struct {
	fileserverHits  atomic.Int32
	db              *database.Queries
	platform        string
	tokenSecret     string
	jwtTTL          time.Duration
	refreshTokenTTL time.Duration
	polkaKey        string
}

func main() {
	cfg := &apiConfig{
		jwtTTL:          time.Hour,
		refreshTokenTTL: 60 * 24 * time.Hour, // 60 days
	}
	err := configureEnv(cfg)
	if err != nil {
		log.Fatal(err)
	}

	serverRootDir := "."
	port := "8080"

	fileserverHanlder := http.StripPrefix(
		"/app",
		http.FileServer(http.Dir(serverRootDir)),
	)

	mux := http.NewServeMux()

	mux.Handle("/app/", cfg.middlewareMetricsInc(fileserverHanlder))

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirpsByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handlerDeleteChirpByID)

	mux.HandleFunc("GET /api/healthz", handlerHealthcheck)

	mux.HandleFunc("POST /api/login", cfg.handlerLogin)

	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)

	mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)

	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", cfg.handlerUpdateUser)

	mux.HandleFunc("POST /api/polka/webhooks", cfg.handlerPolka)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", serverRootDir, port)
	log.Fatal(server.ListenAndServe())
}

func configureEnv(cfg *apiConfig) error {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Printf("DB_URL is required")
		return ErrFailedToLoadEnvVars
	}
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Failed to connect to DB: %s\n", err)
		return ErrDatabaseConnFailed
	}
	dbQueries := database.New(dbConn)

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Printf("PLATFORM is required")
		return ErrFailedToLoadEnvVars
	}

	secret := os.Getenv("TOKEN_SECRET")
	if secret == "" {
		log.Printf("TOKEN_SECRET is required")
		return ErrFailedToLoadEnvVars
	}

	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Printf("POLKA_KEY is required")
		return ErrFailedToLoadEnvVars
	}

	cfg.db = dbQueries
	cfg.platform = platform
	cfg.tokenSecret = secret
	cfg.polkaKey = polkaKey

	return nil
}

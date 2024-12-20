package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/samersawan/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
	polka          string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	polka := os.Getenv("POLKA_KEY")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Could not open connection to the database")
	}
	dbQueries := database.New(db)

	serveMux := http.NewServeMux()
	server := http.Server{Handler: serveMux, Addr: ":8080"}
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		secret:         secret,
		polka:          polka,
	}
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	serveMux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	serveMux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)
	serveMux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	serveMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	serveMux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	serveMux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)
	serveMux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerUpgradeUser)

	server.ListenAndServe()
}

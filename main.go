package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/r4j3sh-com/Goserv/internal/database"
)

type apiConfig struct {
	fileserveHits atomic.Int32
	db            *database.Queries
	platform      string
	tokenSecret   string
	polkaKey      string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	//fmt.Println(dbQueries)
	// Get the platform from the environment
	appPlatform := os.Getenv("PLATFORM")

	jwtSecret := os.Getenv("JWTsecret")
	polkaKey := os.Getenv("POLKA_KEY")

	const port = "8080"
	const filePathRoot = "."

	mux := http.NewServeMux()

	apiCfg := &apiConfig{
		fileserveHits: atomic.Int32{},
		db:            dbQueries,
		platform:      appPlatform,
		tokenSecret:   jwtSecret,
		polkaKey:      polkaKey,
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsers)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUserUpdate)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChipRetriveByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerChirpsDelete)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhooks)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving path %s on port: %s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

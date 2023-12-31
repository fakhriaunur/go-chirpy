package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/fakhriaunur/go-chirpy/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	Secret         string
	PolkaKey       string
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY environment variable is not set")
	}

	const port = "8080"
	const filepathRoot = "."

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if dbg != nil && *dbg {
		if err := db.ResetDB(); err != nil {
			log.Fatal(err)
		}
	}

	apiCfg := &apiConfig{
		fileserverHits: 0,
		DB:             db,
		Secret:         jwtSecret,
		PolkaKey:       polkaKey,
	}

	router := chi.NewRouter()

	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	router.Handle("/app", apiCfg.middlewareMetricsInc(appHandler))
	router.Handle("/app/*", apiCfg.middlewareMetricsInc(appHandler))

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Get("/reset", apiCfg.handlerReset)

	apiRouter.Get("/chirps", apiCfg.handlerChirpsGet)
	apiRouter.Get("/chirps/{chirpID}", apiCfg.handlerChirpsIDGet)
	apiRouter.Post("/chirps", apiCfg.handlerChirpsCreate)
	apiRouter.Delete("/chirps/{chirpID}", apiCfg.handlerChirpsIDDelete)

	apiRouter.Get("/users", apiCfg.handlerUsersGet)
	apiRouter.Post("/users", apiCfg.handlerUsersCreate)

	apiRouter.Put("/users", apiCfg.handlerUserUpdate)
	apiRouter.Post("/login", apiCfg.handlerLogin)
	apiRouter.Post("/refresh", apiCfg.handlerRefresh)
	apiRouter.Post("/revoke", apiCfg.handlerRevoke)

	apiRouter.Post("/polka/webhooks", apiCfg.handlerPolkaWebhooks)

	router.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handlerMetrics)
	router.Mount("/admin", adminRouter)

	corsMux := middlewareCors(router)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

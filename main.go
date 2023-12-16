package main

import (
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
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	const port = "8080"
	const filepathRoot = "."

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := &apiConfig{
		fileserverHits: 0,
		DB:             db,
		Secret:         jwtSecret,
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
	apiRouter.Get("/users", apiCfg.handlerUsersGet)
	apiRouter.Post("/users", apiCfg.handlerUsersCreate)
	apiRouter.Put("/users", apiCfg.handlerUserUpdate)
	apiRouter.Post("/login", apiCfg.handlerLoginCreate)
	apiRouter.Post("/refresh", apiCfg.handlerRefresh)
	apiRouter.Post("/revoke", apiCfg.handlerRevoke)
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

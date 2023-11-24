package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/healthz", http.FileServer(http.Dir(filepathRoot)))
	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const aca = "Access-Control-Allow"
		w.Header().Set(aca+"-Origin", "*")
		w.Header().Set(aca+"-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set(aca+"-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

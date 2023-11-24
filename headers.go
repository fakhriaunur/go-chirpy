package main

import "net/http"

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

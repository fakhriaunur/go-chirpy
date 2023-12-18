package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/fakhriaunur/go-chirpy/internal/auth"
)

func (c *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	log.Println("handlerRevoke")
	refreshToken := r.Header.Get("Authorization")
	refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")

	// log.Println("refreshToken: " + refreshToken)

	issuer, _, err := auth.ValidateToken(refreshToken, c.Secret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't validate token")
		return
	}

	if issuer != "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "issuer is invalid")
		return
	}

	if _, err := c.DB.RevokeSession(refreshToken); err != nil {
		// log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke session")
		return
	}

	// log.Println("success")
	w.WriteHeader(http.StatusOK)
}

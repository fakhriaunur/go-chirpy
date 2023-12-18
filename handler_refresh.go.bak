package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/fakhriaunur/go-chirpy/internal/auth"
)

func (c *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	log.Println("handlerRefresh")
	type returnVals struct {
		Token string `json:"token"`
	}

	refreshToken := r.Header.Get("Authorization")
	refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")

	issuer, userID, err := auth.ValidateToken(refreshToken, c.Secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token is unauthorized")
		return
	}

	if issuer != "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "issuer is unauthorized")
		return
	}

	dbSession, err := c.DB.GetSession(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't retrieve session")
		return
	}

	if dbSession.IsRevoked {
		respondWithError(w, http.StatusUnauthorized, "session is revoked")
		return
	}

	token, err := auth.GenerateToken(c.Secret, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't generate token")
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Token: token,
	})
}

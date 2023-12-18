package main

import (
	"net/http"

	"github.com/fakhriaunur/go-chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	// log.Println("handlerRefresh")
	type returnVals struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't find jwt")
		return
	}

	isRevoked, err := cfg.DB.IsTokenRevoked(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't check session")
		return
	}

	if isRevoked {
		respondWithError(w, http.StatusUnauthorized, "refresh token is revoked")
		return
	}

	accessToken, err := auth.RefreshToken(refreshToken, cfg.Secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't validate jwt")
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Token: accessToken,
	})
}

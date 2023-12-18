package main

import (
	"net/http"

	"github.com/fakhriaunur/go-chirpy/internal/auth"
)

func (c *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	// log.Println("handlerRevoke")
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't find jwt")
		return
	}

	if err := c.DB.RevokeToken(refreshToken); err != nil {
		// log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke session")
		return
	}

	// log.Println("success")
	respondWithJSON(w, http.StatusOK, struct{}{})
}

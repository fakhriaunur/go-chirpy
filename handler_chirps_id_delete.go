package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/fakhriaunur/go-chirpy/internal/auth"
	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerChirpsIDDelete(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	issuer, authorID, err := auth.ValidateToken(token, cfg.Secret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't validate token")
		return
	}

	if issuer != "chirpy-access" {
		respondWithError(w, http.StatusBadRequest, "invalid issuer type")
		return
	}

	chirpID, err := strconv.Atoi(chi.URLParam(r, "chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp id")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't get chirp")
	}

	if dbChirp.AuthorID != authorID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err := cfg.DB.DeleteChirp(chirpID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't delete chirp")
		return
	}

	w.WriteHeader(http.StatusOK)
}

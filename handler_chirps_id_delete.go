package main

import (
	"net/http"
	"strconv"

	"github.com/fakhriaunur/go-chirpy/internal/auth"
	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerChirpsIDDelete(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't get jwt")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't validate token")
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

	authorID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't parse authorID")
		return
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

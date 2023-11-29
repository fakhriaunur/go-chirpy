package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerChirpsIDGet(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(chi.URLParam(r, "chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't get chirp")
		return
	}

	chirp := Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

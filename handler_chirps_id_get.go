package main

import (
	"net/http"
	"strconv"

	"github.com/fakhriaunur/go-chirpy/internal/database"
	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerChirpsIDGet(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		database.Chirp
	}

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

	respondWithJSON(w, http.StatusOK, returnVals{
		Chirp: database.Chirp{
			ID:       dbChirp.ID,
			Body:     dbChirp.Body,
			AuthorID: dbChirp.AuthorID,
		},
	})
}

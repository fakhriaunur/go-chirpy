package main

import (
	"net/http"
	"strconv"

	"github.com/fakhriaunur/go-chirpy/internal/database"
	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerChirpsIDGet(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(chi.URLParam(r, "chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ChirpID is not interger")
		return
	}

	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retreive chirps")
		return
	}

	if chirpID <= 0 || chirpID > len(dbChirps) {
		respondWithError(w, http.StatusNotFound, "Chirp is not found")
		return
	}

	chirp := database.Chirp{}
	for _, dbChirp := range dbChirps {
		if dbChirp.ID == chirpID {
			chirp = dbChirp
		}
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

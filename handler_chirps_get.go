package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/fakhriaunur/go-chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	// log.Println("handlerChirpsGet")
	type returnVals struct {
		Chirps []database.Chirp
	}

	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	hasAuthorIDQuery := r.URL.Query().Has("author_id")
	authorIDStr := r.URL.Query().Get("author_id")
	authorID, err := strconv.Atoi(authorIDStr)
	if err != nil {
		if hasAuthorIDQuery {
			respondWithError(w, http.StatusInternalServerError, "couldn't convert string to int")
			return
		}
	}

	hasSortQuery := r.URL.Query().Has("sort")
	sortQuery := r.URL.Query().Get("sort")

	chirps := []database.Chirp{}
	for _, dbChirp := range dbChirps {
		if hasAuthorIDQuery && dbChirp.AuthorID != authorID {
			continue
		}
		chirps = append(chirps, database.Chirp{
			AuthorID: dbChirp.AuthorID,
			Body:     dbChirp.Body,
			ID:       dbChirp.ID,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		if hasSortQuery && sortQuery == "desc" {
			return chirps[i].ID > chirps[j].ID
		}
		return chirps[i].ID < chirps[j].ID
	})

	// log.Println("success")
	respondWithJSON(w, http.StatusOK, returnVals{
		Chirps: chirps,
	})
}

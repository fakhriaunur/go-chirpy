package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/fakhriaunur/go-chirpy/internal/auth"
	"github.com/fakhriaunur/go-chirpy/internal/database"
)

// type Chirp struct {
// 	ID   int    `json:"id"`
// 	Body string `json:"body"`
// }

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		database.Chirp
	}

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	issuer, authorID, err := auth.ValidateToken(token, cfg.Secret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't validate token")
		return
	}

	if issuer != "chirpy-access" {
		respondWithError(w, http.StatusUnauthorized, "issuer is invalid")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.DB.CreateChirp(authorID, cleaned)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		Chirp: database.Chirp{
			AuthorID: chirp.AuthorID,
			Body:     chirp.Body,
			ID:       chirp.ID,
		},
	})

	// respondWithJSON(w, http.StatusCreated, Chirp{
	// 	ID:   chirp.ID,
	// 	Body: chirp.Body,
	// })
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleanedBody := getCleanedBody(body, badWords)
	return cleanedBody, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	const wordReplacement = "****"
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = wordReplacement
		}
	}
	cleanedBody := strings.Join(words, " ")
	return cleanedBody
}

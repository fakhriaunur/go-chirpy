package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't get jwt")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		// log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couldn't validate token")
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

	authorID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't parse authorID")
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

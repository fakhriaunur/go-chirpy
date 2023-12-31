package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fakhriaunur/go-chirpy/internal/auth"
	"github.com/fakhriaunur/go-chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type returnVals struct {
		database.User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	hashedPwd, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email, hashedPwd)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusConflict, "user already exists")
			return
		}

		respondWithError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		User: database.User{
			Email:       user.Email,
			ID:          user.ID,
			IsChirpyRed: user.IsChirpyRed,
		},
	})

}

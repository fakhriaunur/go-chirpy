package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fakhriaunur/go-chirpy/internal/auth"
	"github.com/fakhriaunur/go-chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type returnVals struct {
		database.User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	dbUser, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't get user")
		return
	}

	if err := auth.CheckHashPassword(dbUser.Password, params.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't validate password")
		return
	}

	token, err := auth.MakeJWT(
		dbUser.ID,
		cfg.Secret,
		time.Hour,
		auth.TokenTypeAccess,
	)
	if err != nil {
		// log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couldn't generate access token")
		return
	}

	refreshToken, err := auth.MakeJWT(
		dbUser.ID,
		cfg.Secret,
		time.Hour*24*30*6,
		auth.TokenTypeRefresh,
	)
	if err != nil {
		// log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couldn't generate refresh token")
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		User: database.User{
			ID:          dbUser.ID,
			Email:       dbUser.Email,
			IsChirpyRed: dbUser.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}

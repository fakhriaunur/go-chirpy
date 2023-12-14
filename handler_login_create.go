package main

import (
	"encoding/json"
	"net/http"

	"github.com/fakhriaunur/go-chirpy/internal/auth"
)

func (c *apiConfig) handlerLoginCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Expiry   *int   `json:"expires_in_seconds"`
	}

	type returnVals struct {
		Email string `json:"email"`
		ID    int    `json:"id"`
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	dbUser, err := c.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't get user")
		return
	}

	if err := auth.CheckHashPassword(dbUser.Password, params.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't validate password")
		return
	}

	token, err := auth.GenerateToken(c.Secret, dbUser.ID, params.Expiry)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't generate token")
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Email: dbUser.Email,
		ID:    dbUser.ID,
		Token: token,
	})

}

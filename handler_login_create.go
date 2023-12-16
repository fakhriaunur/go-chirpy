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
	}

	type returnVals struct {
		ID           int    `json:"id"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	// dbToken, err := c. DB.

	if err := auth.CheckHashPassword(dbUser.Password, params.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't validate password")
		return
	}

	token, err := auth.GenerateToken(c.Secret, dbUser.ID)
	if err != nil {
		// log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couldn't generate token")
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(c.Secret, dbUser.ID)
	if err != nil {
		// log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couldn't generate refresh token")
		return
	}

	// log.Println("create session")
	if _, err := c.DB.CreateSession(refreshToken); err != nil {
		// log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couldn't create session")
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		ID:           dbUser.ID,
		Email:        dbUser.Email,
		Token:        token,
		RefreshToken: refreshToken,
	})
}

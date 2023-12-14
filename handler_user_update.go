package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/fakhriaunur/go-chirpy/internal/auth"
)

func (c *apiConfig) handlerUserUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type returnVals struct {
		Email string `json:"email"`
		ID    int    `json:"id"`
	}

	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	// log.Println(token)

	userID, err := auth.ValidateToken(token, c.Secret)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, "couldn't validate token")
		return
	}
	// log.Println(userID)

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		// log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	newPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password")
		return
	}
	// log.Println(newPassword)

	dbUser, err := c.DB.UpdateUser(userID, params.Email, newPassword)
	if err != nil {
		respondWithError(w, http.StatusConflict, "couldnt update user")
		return
	}

	// log.Println(dbUser)

	respondWithJSON(w, http.StatusOK, returnVals{
		Email: dbUser.Email,
		ID:    dbUser.ID,
	})
}

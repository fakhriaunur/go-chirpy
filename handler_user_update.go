package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fakhriaunur/go-chirpy/internal/auth"
	"github.com/fakhriaunur/go-chirpy/internal/database"
)

func (cfg *apiConfig) handlerUserUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type returnVals struct {
		database.User
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't get jwt")
		return
	}
	// log.Println(token)

	subject, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		// log.Println(err)
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
	userID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't parse userID")
		return
	}

	dbUser, err := cfg.DB.UpdateUser(userID, params.Email, newPassword)
	if err != nil {
		respondWithError(w, http.StatusConflict, "couldnt update user")
		return
	}

	// log.Println(dbUser)

	respondWithJSON(w, http.StatusOK, returnVals{
		User: database.User{
			Email:       dbUser.Email,
			ID:          dbUser.ID,
			IsChirpyRed: dbUser.IsChirpyRed,
		},
	})
}

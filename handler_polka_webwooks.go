package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fakhriaunur/go-chirpy/internal/auth"
	"github.com/fakhriaunur/go-chirpy/internal/database"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		}
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't find api key")
		return
	}
	if apiKey != cfg.PolkaKey {
		respondWithError(w, http.StatusUnauthorized, "API Key is invalid")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, struct{}{})
		return
	}

	if _, err := cfg.DB.UpgradeUser(params.Data.UserID); err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "couldn't find user")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "couldn't upgrade user")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}

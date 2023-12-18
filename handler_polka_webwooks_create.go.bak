package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerPolkaWebhooksCreate(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		UserID int `json:"user_id"`
	}

	type parameters struct {
		Event string `json:"event"`
		Data  `json:"data"`
	}

	polkaKey := r.Header.Get("Authorization")
	polkaKey = strings.TrimPrefix(polkaKey, "ApiKey ")

	if polkaKey != cfg.PolkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusOK)
		return
	}

	dbUser, err := cfg.DB.GetUser(params.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if _, err := cfg.DB.UpgradeUser(dbUser.ID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't upgrade user")
		return
	}

	// w.WriteHeader(http.StatusOK)
}

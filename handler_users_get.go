package main

import (
	"net/http"
	"sort"

	"github.com/fakhriaunur/go-chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request) {
	dbUsers, err := cfg.DB.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't retrieve users")
	}

	users := []database.User{}
	for _, dbUser := range dbUsers {
		users = append(users, database.User{
			ID:    dbUser.ID,
			Email: dbUser.Email,
		})
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	respondWithJSON(w, http.StatusOK, users)
}

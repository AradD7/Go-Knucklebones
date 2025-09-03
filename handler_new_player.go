package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AradD7/Go-Knuclebones/internal/auth"
	"github.com/AradD7/Go-Knuclebones/internal/database"
	"github.com/google/uuid"
)

type Player struct {
	Id 			uuid.UUID `json:"id"`
	CreatedAt	time.Time `json:"created_at"`
	Username 	string	  `json:"username"`
}

func (cfg *apiConfig) handlerNewPlayer(w http.ResponseWriter, r *http.Request) {
	type createPlayer struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	var newPlayer createPlayer
	if err := decoder.Decode(&newPlayer); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode the json data", err)
		return
	}

	hashPassword, err := auth.HashPassword(newPlayer.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to has the password", err)
		return
	}

	player, err := cfg.db.CreatePlayer(r.Context(), database.CreatePlayerParams{
		Username: 		newPlayer.Username,
		HashedPassword: hashPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to add the player to DB", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Player{
		Id: 		player.ID,
		CreatedAt: 	player.CreatedAt,
		Username: 	player.Username,
	})
}

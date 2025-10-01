package main

import (
	"math/rand"
	"net/http"

	"github.com/AradD7/Go-Knuclebones/internal/auth"
	"github.com/google/uuid"
)

type DiceRoll struct {
	Dice int `json:"dice"`
}

func (cfg *apiConfig) handlerRoll(w http.ResponseWriter, r *http.Request) {
	gameId, err := uuid.Parse(r.URL.Query().Get("game_id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Game ID is not valid", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not Authorized", err)
		return
	}

	_, err = auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token is exipred, refresh JWT token or login again", err)
		return
	}

	dice := rand.Intn(6) + 1

	cfg.gs.broadcastRolled(gameId, dice)

	respondWithJSON(w, http.StatusOK, DiceRoll{
		Dice: dice,
	})
}

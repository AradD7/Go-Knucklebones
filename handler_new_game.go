package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AradD7/Go-Knuclebones/internal/auth"
	"github.com/AradD7/Go-Knuclebones/internal/database"
	"github.com/google/uuid"
)

type InitGame struct{
	Id 					uuid.UUID `json:"id"`
	CreatedAt 			time.Time `json:"created_at"`
	Player1Username 	string 	  `json:"player1_username"`
	Player2Username 	string 	  `json:"player2_username"`
}

func (cfg *apiConfig) handlerNewGame(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to get JWT token from request header", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var params parameters
	if err = decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode json data", err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token is exipred, refresh JWT token or login again", err)
		return
	}

	player1, err := cfg.db.GetPlayerByPlayerId(r.Context(), userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Faild to get player from DB", err)
		return
	}

	player2, err := cfg.db.GetPlayerByUsername(r.Context(), params.Username)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Faild to get player2 from DB", err)
		return
	}

	player1Board, err := cfg.db.CreateBoard(r.Context(), player1.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Faild to initialize board for player1", err)
		return
	}

	player2Board, err := cfg.db.CreateBoard(r.Context(), player2.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Faild to initialize board for player2", err)
		return
	}

	newGame, err := cfg.db.CreateNewGame(r.Context(), database.CreateNewGameParams{
		Board1: player1Board.ID,
		Board2: player2Board.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Faild to create a game", err)
		return
	}

	if err = cfg.db.LinkGame(r.Context(), database.LinkGameParams{
		GameID: uuid.NullUUID{
			Valid: true,
			UUID:  newGame.ID,
		},
		ID: 	player1Board.ID,
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Faild to link board1 to game", err)
		return
	}
	if err = cfg.db.LinkGame(r.Context(), database.LinkGameParams{
		GameID: uuid.NullUUID{
			Valid: true,
			UUID:  newGame.ID,
		},
		ID: 	player2Board.ID,
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Faild to link board2 to game", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, InitGame{
		Id: 				newGame.ID,
		CreatedAt: 			newGame.CreatedAt,
		Player1Username:  	player1.Username,
		Player2Username: 	player2.Username,
	})
}

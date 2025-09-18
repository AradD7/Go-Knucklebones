package main

import (
	"net/http"
	"time"

	"github.com/AradD7/Go-Knuclebones/internal/auth"
	"github.com/AradD7/Go-Knuclebones/internal/database"
	"github.com/google/uuid"
)

type Game struct{
	Id 			uuid.UUID `json:"id"`
	CreatedAt	time.Time `json:"created_at"`
	Board1 		[][]int32 `json:"board1"`
	Board2		[][]int32 `json:"board2"`
}

type GameIds struct {
	Ids 	[]uuid.UUID `json:"ids"`
}

func (cfg *apiConfig) handlerNewGame(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not Authorized", err)
		return
	}

	playerId, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token is exipred, refresh JWT token or login again", err)
		return
	}

	player1, err := cfg.db.GetPlayerByPlayerId(r.Context(), playerId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Faild to get player from DB", err)
		return
	}

	player1Board, err := cfg.db.CreateBoard(r.Context(), player1.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Faild to initialize board for player1", err)
		return
	}

	newGame, err := cfg.db.CreateNewGame(r.Context(), database.CreateNewGameParams{
		Board1: player1Board.ID,
		Board2: uuid.NullUUID{Valid: false},
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

	respondWithJSON(w, http.StatusCreated, Game{
		Id: 		newGame.ID,
		CreatedAt: 	newGame.CreatedAt,
		Board1: 	player1Board.Board,
		Board2: 	nil,
	})
}

func (cfg *apiConfig) handlerGetGames(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to get JWT token from request header", err)
		return
	}

	playerId, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token is exipred, refresh JWT token or login again", err)
		return
	}

	gameIds, err := cfg.db.GetGamesWithPlayerId(r.Context(), playerId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Faild to get player games from DB", err)
		return
	}

	var games GameIds
	for _, id := range gameIds {
		games.Ids = append(games.Ids, id.UUID)
	}

	respondWithJSON(w, http.StatusOK, games)
}

func (cfg *apiConfig) handlerGetGame(w http.ResponseWriter, r *http.Request) {
	gameId, err := uuid.Parse(r.PathValue("game_id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Game ID is not valid", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not Authorized", err)
		return
	}

	playerId, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token is exipred, refresh JWT token or login again", err)
		return
	}

	game, err := cfg.db.GetGameById(r.Context(), gameId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Faild to get game from DB", err)
		return
	}

	board1, err := cfg.db.GetBoardById(r.Context(), game.Board1)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to get board 1 of the game", err)
		return
	}

	board2, err := cfg.db.GetBoardById(r.Context(), game.Board2)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to get board 2 of the game", err)
		return
	}

	if board1.PlayerID == playerId {
		respondWithJSON(w, http.StatusOK, Game{
			Id: 		game.ID,
			CreatedAt: 	game.CreatedAt,
			Board1:  	board1.Board,
			Board2: 	board2.Board,
		})
		return
	}

	if board2.PlayerID == playerId {
		respondWithJSON(w, http.StatusOK, Game{
			Id: 		game.ID,
			CreatedAt: 	game.CreatedAt,
			Board1:  	board2.Board,
			Board2: 	board1.Board,
		})
		return
	}
	respondWithError(w, http.StatusUnauthorized, "Player is not in this game", err)
	return
}

func (cfg *apiConfig) handlerJoinGame(w http.ResponseWriter, r *http.Request) {
	gameId, err := uuid.Parse(r.PathValue("game_id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Game ID is not valid", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not Authorized", err)
		return
	}

	playerId, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token is exipred, refresh JWT token or login again", err)
		return
	}

	currentGame, err := cfg.db.GetGameById(r.Context(), gameId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Game not found", err)
		return
	}

	playerBoard, err := cfg.db.CreateBoard(r.Context(), playerId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Falied to initialize board", err)
		return
	}

	if err = cfg.db.JoinGame(r.Context(), database.JoinGameParams{
		Board2: uuid.NullUUID{
			Valid: true,
			UUID:  playerBoard.ID,
		},
		ID: 	gameId,
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Falied to join game", err)
		return
	}

	if err = cfg.db.LinkGame(r.Context(), database.LinkGameParams{
		ID: 	playerBoard.ID,
		GameID: uuid.NullUUID{
			Valid: true,
			UUID:  gameId,
		},
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Falied to link game to board", err)
		return
	}

	oppBoard, err := cfg.db.GetBoardById(r.Context(), currentGame.Board1)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Falied to get opponent board", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Game{
		Id: 		gameId,
		CreatedAt: 	currentGame.CreatedAt,
		Board1: 	playerBoard.Board,
		Board2: 	oppBoard.Board,
	})

	cfg.gs.broadcastToGame(gameId)
}

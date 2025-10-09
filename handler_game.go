package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/AradD7/Go-Knuclebones/internal/auth"
	"github.com/AradD7/Go-Knuclebones/internal/database"
	"github.com/google/uuid"
)

type Game struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Board1    [][]int32 `json:"board1"`
	Board2    [][]int32 `json:"board2"`
	Score1    int       `json:"score1"`
	Score2    int       `json:"score2"`
	IsTurn    bool      `json:"is_turn"`
	IsOver    bool      `json:"is_over"`
	OppName   string    `json:"opp_name"`
	OppAvatar string    `json:"opp_avatar"`
}

type GameOverview struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"date"`
	Status    int	    `json:"status"` //-1 means lost, 0 means in progress, 1 means won
	OppName   string    `json:"opp_name"`
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

	_ = cfg.db.DeleteEmptyBoardsForPlayer(r.Context(), playerId) //don't care if error, this is just housekeeping

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
		ID: player1Board.ID,
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Faild to link board1 to game", err)
		return
	}

	var player1BoardData [][]int32
	if err = json.Unmarshal(player1Board.Board, &player1BoardData); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't turn the board into [][]int32", err)
		return
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	respondWithJSON(w, http.StatusCreated, Game{
		Id:        newGame.ID,
		CreatedAt: newGame.CreatedAt,
		Board1:    player1BoardData,
		Board2:    nil,
	})
}

func (cfg *apiConfig) handlerGetGames(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to get JWT token from request header", err)
		return
	}

	playerId, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token is exipred, refresh JWT token or login again", err)
		return
	}

	playerGames, err := cfg.db.GetGamesWithPlayerId(r.Context(), playerId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Faild to get player games from DB", err)
		return
	}

	var games []GameOverview
	for _, game := range playerGames {
		status := 0
		if game.WinnerID.Valid {
			switch game.WinnerID.UUID {
			case playerId:
				status = 1
			default:
				status = -1
			}
		}
		games = append(games, GameOverview{
			Id: 		game.GameID,
			CreatedAt: 	game.Date,
			OppName: 	game.OpponentName,
			Status: 	status,
		})
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
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

	board2, err := cfg.db.GetBoardById(r.Context(), game.Board2.UUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to get board 2 of the game", err)
		return
	}

	var board1Data, board2Data [][]int32
	if err = json.Unmarshal(board1.Board, &board1Data); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't turn the board into [][]int32", err)
		return
	}
	if err = json.Unmarshal(board2.Board, &board2Data); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't turn the board into [][]int32", err)
		return
	}


	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	if board1.PlayerID == playerId {
		opp, err := cfg.db.GetPlayerByPlayerId(r.Context(), board2.PlayerID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to get player info to broadcast", err)
			return
		}

		oppDisplayName := opp.Username
		if opp.DisplayName.Valid {
			oppDisplayName = opp.DisplayName.String
		}

		respondWithJSON(w, http.StatusOK, Game{
			Id:        game.ID,
			CreatedAt: game.CreatedAt,
			Board1:    board1Data,
			Board2:    board2Data,
			Score1:    int(board1.Score.Int32),
			Score2:    int(board2.Score.Int32),
			OppName:   oppDisplayName,
			OppAvatar: opp.Avatar.String,
			IsTurn:    game.PlayerTurn.UUID == playerId,
			IsOver:    game.Winner.Valid,
		})
		return
	}

	if board2.PlayerID == playerId {
		opp, err := cfg.db.GetPlayerByPlayerId(r.Context(), board1.PlayerID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "failed to get player info to broadcast", err)
			return
		}

		oppDisplayName := opp.Username
		if opp.DisplayName.Valid {
			oppDisplayName = opp.DisplayName.String
		}

		respondWithJSON(w, http.StatusOK, Game{
			Id:        game.ID,
			CreatedAt: game.CreatedAt,
			Board1:    board2Data,
			Board2:    board1Data,
			Score1:    int(board2.Score.Int32),
			Score2:    int(board1.Score.Int32),
			OppName:   oppDisplayName,
			OppAvatar: opp.Avatar.String,
			IsTurn:    game.PlayerTurn.UUID == playerId,
			IsOver:    game.Winner.Valid,
		})
		return
	}
	respondWithError(w, http.StatusUnauthorized, "Player is not in this game", err)
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

	if currentGame.Board2.Valid {
		respondWithError(w, http.StatusConflict, "Already in game", nil)
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
		ID: gameId,
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Falied to join game", err)
		return
	}

	if err = cfg.db.LinkGame(r.Context(), database.LinkGameParams{
		ID: playerBoard.ID,
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

	var playerBoardData, oppBoardData [][]int32
	if err = json.Unmarshal(playerBoard.Board, &playerBoardData); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't turn the board into [][]int32", err)
		return
	}
	if err = json.Unmarshal(oppBoard.Board, &oppBoardData); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't turn the board into [][]int32", err)
		return
	}

	var playerTurnId uuid.UUID
	if rand.Intn(2) == 0 {
		playerTurnId = oppBoard.PlayerID
	} else {
		playerTurnId = playerId
	}
	if err = cfg.db.SetPlayerTurn(r.Context(), database.SetPlayerTurnParams{
		ID: gameId,
		PlayerTurn: uuid.NullUUID{
			Valid: true,
			UUID:  playerTurnId,
		},
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to assign turn", err)
		return
	}

	opp, err := cfg.db.GetPlayerByPlayerId(r.Context(), oppBoard.PlayerID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get player info to broadcast", err)
		return
	}

	oppDisplayName := opp.Username
	if opp.DisplayName.Valid {
		oppDisplayName = opp.DisplayName.String
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	respondWithJSON(w, http.StatusCreated, Game{
		Id:        gameId,
		CreatedAt: currentGame.CreatedAt,
		Board1:    playerBoardData,
		Board2:    oppBoardData,
		IsTurn:    playerId == playerTurnId,
		OppName:   oppDisplayName,
		OppAvatar: opp.Avatar.String,
	})

	player, err := cfg.db.GetPlayerByPlayerId(r.Context(), playerId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get player info to broadcast", err)
		return
	}

	if player.DisplayName.Valid {
		cfg.gs.broadcastJoined(gameId, player.DisplayName.String, player.Avatar.String)
	} else {
		cfg.gs.broadcastJoined(gameId, player.Username, player.Avatar.String)
	}
}

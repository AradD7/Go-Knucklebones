package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AradD7/Go-Knuclebones/internal/auth"
	"github.com/AradD7/Go-Knuclebones/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerMakeMove(w http.ResponseWriter, r *http.Request) {
	gameId, err := uuid.Parse(r.PathValue("game_id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Game ID is not valid", err)
		return
	}

	type moveParameters struct {
		Dice  int `json:"dice"`
		Row   int `json:"row"`
		Col   int `json:"col"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Not Authorized", err)
		return
	}

	playerId, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token is exipred, refresh JWT token or login again", err)
		return
	}

	var move moveParameters
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&move); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode json data", err)
		return
	}

	playerBoard, err := cfg.db.GetBoardByPlayerIdAndGameId(r.Context(), database.GetBoardByPlayerIdAndGameIdParams{
		PlayerID: 	playerId,
		GameID: 	uuid.NullUUID{
			Valid: 	true,
			UUID: 	gameId,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Player is not in this game", err)
		return
	}

	currentGame, err := cfg.db.GetGameById(r.Context(), gameId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Faild to get game from DB", err)
		return
	}

	oppBoardId := currentGame.Board1
	if oppBoardId == playerBoard.ID {
		oppBoardId = currentGame.Board2.UUID
	}

	oppBoard, err := cfg.db.GetBoardById(r.Context(), oppBoardId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Opponent not found", err)
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

	updatedPlayerBoard, err := putDice(playerBoardData, move.Dice, move.Row, move.Col)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Can't put there!", err)
		return
	}

	updatedOppBoard := updateOpp(oppBoardData, move.Dice, move.Col)

	if isFull(updatedPlayerBoard) {
		if calcScore(updatedPlayerBoard) > calcScore(updatedOppBoard) {
			cfg.db.SetGameWinner(r.Context(), database.SetGameWinnerParams{
				ID: 	currentGame.ID,
				Winner: uuid.NullUUID{
					Valid: true,
					UUID:  playerBoard.PlayerID,
				},
			})
		} else {
			cfg.db.SetGameWinner(r.Context(), database.SetGameWinnerParams{
				ID: 	currentGame.ID,
				Winner: uuid.NullUUID{
					Valid: true,
					UUID:  oppBoard.PlayerID,
				},
			})
		}
	}

	updatedPlayerBoardJSON, err := json.Marshal(updatedPlayerBoard)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't turn the board into json data", err)
		return
	}
	updatedOppBoardJSON, err := json.Marshal(updatedOppBoard)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't turn the board into json data", err)
		return
	}

	if err = cfg.db.UpdateBoard(r.Context(), database.UpdateBoardParams{
		ID: playerBoard.ID,
		Board: updatedPlayerBoardJSON,
		Score: sql.NullInt32{
			Valid: true,
			Int32: calcScore(updatedPlayerBoard),
		},
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update player board", err)
		return
	}
	if err = cfg.db.UpdateBoard(r.Context(), database.UpdateBoardParams{
		ID: oppBoardId,
		Board: updatedOppBoardJSON,
		Score: sql.NullInt32{
			Valid: true,
			Int32: calcScore(updatedOppBoard),
		},
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update opponent board", err)
		return
	}

	if err = cfg.db.SetPlayerTurn(r.Context(), database.SetPlayerTurnParams{
		ID: 		currentGame.ID,
		PlayerTurn: uuid.NullUUID{
			Valid: 	true,
			UUID: 	oppBoard.PlayerID,
		},
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to assign turn", err)
		return
	}

	cfg.gs.broadcastToGame(currentGame.ID)

	respondWithJSON(w, http.StatusOK, GameState{
		Board1: updatedPlayerBoard,
		Board2: updatedOppBoard,
		Score1: int(playerBoard.Score.Int32),
		Score2: int(oppBoard.Score.Int32),
		IsOver: isFull(updatedPlayerBoard) || isFull(updatedOppBoard),
	})
}

func calcScore(board [][]int32) int32 {
	score := 0
	for col := range 3 {
		multiplier := make(map[int]int)
		for row := range 3 {
			current_number := int(board[row][col])
			multiplier[current_number] += 1
		}
		for key, val := range multiplier {
			score += (key * val) * val
		}
	}

	return int32(score)
}

func putDice(board [][]int32, dice, row, col int) ([][]int32, error) {
	if row < 0 || row > 3 || col < 0 || col > 3 {
		return board, fmt.Errorf("There are 3 rows and 3 columns (ie 3 > row, col > 0)")
	}

	if board[row][col] != 0 {
		return board, fmt.Errorf("Already full! Place dice in another cell")
	}

	if row < 2 && board[row + 1][col] == 0 {
		return board, fmt.Errorf("Can't place here! Bottom cell is empty")
	}

	board[row][col] = int32(dice)
	return board, nil
}

func updateOpp(board [][]int32, dice, col int) ([][]int32) {
	for row := range 3 {
		if board[row][col] == int32(dice) {
			board[row][col] = 0
		}
	}

	for row := 2; row > 0; row-- {
		if board[row][col] == 0 {
			check := row - 1
			for check >= 0 && board[check][col] == 0 {
				check -= 1
			}
			if check >= 0 {
				board[row][col], board[check][col] = board[check][col], board[row][col]
			}
		}
	}

	return board
}

func isFull(board [][]int32) bool {
	for row := range 3 {
		for col := range 3 {
			if board[row][col] == 0 {
				return false
			}
		}
	}
	return true
}

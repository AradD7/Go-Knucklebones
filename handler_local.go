package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type GameState struct {
	Board1 	[][]int32 	`json:"board1"`
	Board2 	[][]int32 	`json:"board2"`
	Score1 	int			`json:"score1"`
	Score2 	int			`json:"score2"`
	IsOver 	bool 		`json:"is_over"`
}

func (cfg *apiConfig) handlerLocalGame(w http.ResponseWriter, r *http.Request) {
	log.Println("Request recieved!")
	type parameters struct {
		Board1 	[][]int32 	`json:"board1"`
		Board2 	[][]int32 	`json:"board2"`
		Turn 	string 		`json:"turn"` //either "player1" or "player2"
		Dice 	int 		`json:"dice"`
		Row  	int			`json:"row"`
		Col  	int			`json:"col"`
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to read json data", err)
		return
	}

	var updatedGameState GameState
	switch params.Turn {
	case "player1":
		updatedBoard1, err := putDice(params.Board1, params.Dice, params.Row, params.Col)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Can't put dice there", err)
			return
		}
		updatedGameState.Board1 = updatedBoard1
		updatedGameState.Board2 = updateOpp(params.Board2, params.Dice, params.Col)
		updatedGameState.Score1 = int(calcScore(updatedGameState.Board1))
		updatedGameState.Score2 = int(calcScore(updatedGameState.Board2))
		updatedGameState.IsOver = isFull(updatedBoard1)
	case "player2":
		updatedBoard2, err := putDice(params.Board2, params.Dice, params.Row, params.Col)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Can't put dice there", err)
			return
		}
		updatedGameState.Board2 = updatedBoard2
		updatedGameState.Board1 = updateOpp(params.Board1, params.Dice, params.Col)
		updatedGameState.Score2 = int(calcScore(updatedGameState.Board2))
		updatedGameState.Score1 = int(calcScore(updatedGameState.Board1))
		updatedGameState.IsOver = isFull(updatedBoard2)
	default:
		respondWithError(w, http.StatusBadRequest, "Turn field cannot be empty", nil)
		}

	respondWithJSON(w, http.StatusOK, updatedGameState)
}

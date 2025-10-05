package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sort"
)

type gameVsComputer struct {
	Board1     [][]int32 `json:"board1"`
	Board2     [][]int32 `json:"board2"`
	NextBoard1 [][]int32 `json:"next_board1"`
	NextBoard2 [][]int32 `json:"next_board2"`
	Score1     int       `json:"score1"`
	Score2     int       `json:"score2"`
	NextScore1 int       `json:"next_score1"`
	NextScore2 int       `json:"next_score2"`
	NextDice   int       `json:"next_dice"`
	IsOver     bool      `json:"is_over"`
	IsOverNext bool      `json:"is_over_next"`
}

func (cfg *apiConfig) handlerComputerGame(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Board1     [][]int32 `json:"board1"`
		Board2     [][]int32 `json:"board2"`
		Dice       int       `json:"dice"`
		Row        int       `json:"row"`
		Col        int       `json:"col"`
		Difficulty string    `json:"difficulty"`
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to read json data", err)
		return
	}

	var updatedGameState gameVsComputer
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

	if updatedGameState.IsOver {
		respondWithJSON(w, http.StatusOK, updatedGameState)
		return
	}

	nextDice := rand.Intn(6) + 1
	nextboard2, nextboard1 := computerMove(updatedGameState.Board2, updatedGameState.Board1, params.Difficulty, nextDice)
	updatedGameState.NextBoard1 = nextboard1
	updatedGameState.NextBoard2 = nextboard2
	updatedGameState.NextScore1 = int(calcScore(updatedGameState.NextBoard1))
	updatedGameState.NextScore2 = int(calcScore(updatedGameState.NextBoard2))
	updatedGameState.IsOverNext = isFull(updatedGameState.NextBoard2)
	updatedGameState.NextDice = nextDice
	respondWithJSON(w, http.StatusOK, updatedGameState)
}

func computerMove(board1, board2 [][]int32, difficulty string, dice int) ([][]int32, [][]int32) {
	computerBoard := deepCopy2D(board1)
	playerBoard := deepCopy2D(board2)
	type scenario struct {
		board1    [][]int32
		board2    [][]int32
		diffScore int
		filledRow int
	}

	var scenarios []scenario

	for col := range 3 {
		var row int
		if computerBoard[2][col] == 0 {
			row = 2
		} else if computerBoard[1][col] == 0 {
			row = 1
		} else if computerBoard[0][col] == 0 {
			row = 0
		} else {
			continue
		}
		updatedBoard1, _ := putDice(computerBoard, dice, row, col)
		updatedBoard2 := updateOpp(playerBoard, dice, col)
		scenarios = append(scenarios, scenario{
			board1:    updatedBoard1,
			board2:    updatedBoard2,
			diffScore: int(calcScore(updatedBoard1)) - int(calcScore(updatedBoard2)),
			filledRow: row,
		})
	}

	sort.Slice(scenarios, func(i, j int) bool {
		if scenarios[i].diffScore != scenarios[j].diffScore {
			return scenarios[i].diffScore > scenarios[j].diffScore
		}
		return scenarios[i].filledRow >= scenarios[j].filledRow
	})

	scenarioIdx := 0
	switch len(scenarios) {
	case 3:
		switch difficulty {
		case "easy":
			scenarioIdx = 2
		case "medium":
			scenarioIdx = 1
		case "hard":
			scenarioIdx = 0
		}
	case 2:
		if difficulty == "hard" {
			scenarioIdx = 0
		} else {
			scenarioIdx = 1
		}
	default:
		scenarioIdx = 0
	}
	return scenarios[scenarioIdx].board1, scenarios[scenarioIdx].board2
}

func deepCopy2D(original [][]int32) [][]int32 {
	copied := make([][]int32, len(original))
	for i := range original {
		copied[i] = make([]int32, len(original[i]))
		copy(copied[i], original[i])
	}
	return copied
}

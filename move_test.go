package main

import "testing"

func TestMove(t *testing.T) {
	type TestTable struct {
		playerBoard 		[][][]int32
		diceRowCol 			[][]int
		oppBoard 			[][][]int32
		updatedPlayerBoard 	[][][]int32
		updatedOppBoard 	[][][]int32
		playerScore 		[]int32
	}

	tests := TestTable{
		playerBoard: [][][]int32{
			{
				{0, 1, 0},
				{0, 3, 5},
				{3, 2, 6},
			},
			{
				{1, 0, 0},
				{3, 4, 0},
				{6, 1, 3},
			},
			{
				{4, 0, 2},
				{1, 6, 2},
				{1, 4, 2},
			},
			{
				{6, 3, 0},
				{6, 3, 0},
				{6, 3, 0},
			},
			{
				{6, 2, 0},
				{1, 3, 0},
				{6, 3, 0},
			},
		},
		diceRowCol: [][]int{
			{6, 1, 0},
			{5, 0, 1},
			{3, 0, 1},
			{1, 2, 2},
			{3, 2, 2},
		},
		oppBoard: [][][]int32{
			{
				{6, 0, 0},
				{4, 0, 0},
				{6, 0, 0},
			},
			{
				{0, 0, 0},
				{0, 5, 0},
				{0, 5, 0},
			},
			{
				{0, 0, 0},
				{0, 0, 0},
				{0, 3, 0},
			},
			{
				{0, 0, 3},
				{0, 0, 1},
				{0, 0, 1},
			},
			{
				{0, 0, 3},
				{0, 0, 3},
				{0, 0, 1},
			},
		},
		updatedPlayerBoard: [][][]int32{
			{
				{0, 1, 0},
				{6, 3, 5},
				{3, 2, 6},
			},
			{
				{1, 5, 0},
				{3, 4, 0},
				{6, 1, 3},
			},
			{
				{4, 3, 2},
				{1, 6, 2},
				{1, 4, 2},
			},
			{
				{6, 3, 0},
				{6, 3, 0},
				{6, 3, 1},
			},
			{
				{6, 2, 0},
				{1, 3, 0},
				{6, 3, 3},
			},
		},
		updatedOppBoard: [][][]int32{
			{
				{0, 0, 0},
				{0, 0, 0},
				{4, 0, 0},
			},
			{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 3},
			},
			{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 1},
			},
		},
		playerScore: []int32{
			26,
			23,
			39,
			82,
			42,
		},
	}

	for i := range 5 {
		getPlayerBoard, err := putDice(tests.playerBoard[i], tests.diceRowCol[i][0], tests.diceRowCol[i][1], tests.diceRowCol[i][2])
		if err != nil {
			t.Fatalf("Shouldn't have gotten an error: %v", err)
		}
		getOppBoard := updateOpp(tests.oppBoard[i], tests.diceRowCol[i][0], tests.diceRowCol[i][2])

		if !testEqual(getPlayerBoard, tests.updatedPlayerBoard[i]) {
			t.Fatalf("Expected %v to equal %v", getPlayerBoard, tests.updatedPlayerBoard[i])
		}
		if !testEqual(getOppBoard, tests.updatedOppBoard[i]) {
			t.Fatalf("Expected %v to equal %v", getOppBoard, tests.updatedOppBoard[i])
		}
		if calcScore(getPlayerBoard) != tests.playerScore[i] {
			t.Fatalf("Got score %d, wanted score %d", calcScore(getPlayerBoard), tests.playerScore[i])
		}
	}
}

func testEqual(board1, board2 [][]int32) bool {
	for i, row := range board1 {
		for j, val := range row {
			if val != board2[i][j] {
				return false
			}
		}
	}
	return true
}

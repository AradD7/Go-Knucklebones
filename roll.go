package main

import (
	"math/rand"
	"net/http"
)

func handlerRoll(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, struct{
		Dice int `json:"dice"`
	}{
			Dice: rand.Intn(6) + 1,
		})
}

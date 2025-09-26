package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerComputerGame(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		IdToken string `json:"id_token"`
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to read json data", err)
		return
	}


}

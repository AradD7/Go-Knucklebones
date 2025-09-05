package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusUnauthorized, "Reset is only possible in dev mode", nil)
		return
	}

	cfg.db.ResetDatabase(r.Context())
	respondWithJSON(w, http.StatusNoContent, nil)
	log.Println("Successfully emptied the database!")
}

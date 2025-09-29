package main

import "net/http"

func (cfg *apiConfig) handlerHealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, struct{
		Status string `json:"status"`
	}{
			Status: "ok",
		})
}

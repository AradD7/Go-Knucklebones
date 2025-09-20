package main

import (
	"net/http"
	"time"

	"github.com/AradD7/Go-Knuclebones/internal/auth"
)

type newToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No token found in the header", err)
		return
	}

	refreshToken, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil || time.Now().After(refreshToken.ExpiresAt) || refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token not valid. Login again", err)
		return
	}

	jwtToken, err := auth.MakeJWT(refreshToken.PlayerID, cfg.tokenSecret, time.Minute * 60)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create new JWT token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, newToken{
		Token: jwtToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No token found in the header", err)
		return
	}

	_ = cfg.db.RevokeRefreshToken(r.Context(), token)

	respondWithJSON(w, http.StatusOK, nil)
}

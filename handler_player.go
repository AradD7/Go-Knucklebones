package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AradD7/Go-Knuclebones/internal/auth"
	"github.com/AradD7/Go-Knuclebones/internal/database"
	"github.com/google/uuid"
)

type Player struct {
	Id 			 uuid.UUID 	`json:"id"`
	CreatedAt	 time.Time 	`json:"created_at"`
	Username 	 string	   	`json:"username"`
	RefreshToken string 	`json:"refresh_token"`
	Token 		 string 	`json:"token"`
	Avatar 		 string 	`json:"avatar"`
}

func (cfg *apiConfig) handlerNewPlayer(w http.ResponseWriter, r *http.Request) {
	type createPlayerParams struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	var newPlayer createPlayerParams
	if err := decoder.Decode(&newPlayer); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode the json data", err)
		return
	}

	hashPassword, err := auth.HashPassword(newPlayer.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to has the password", err)
		return
	}

	player, err := cfg.db.CreatePlayer(r.Context(), database.CreatePlayerParams{
		Username: 		newPlayer.Username,
		HashedPassword: hashPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to add the player to DB", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Player{
		Id: 		player.ID,
		CreatedAt: 	player.CreatedAt,
		Username: 	player.Username,
	})
}

func (cfg *apiConfig) handlerPlayerLogin(w http.ResponseWriter, r *http.Request) {
	type loginPlayerParams struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	var newPlayer loginPlayerParams
	if err := decoder.Decode(&newPlayer); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode the json data", err)
		return
	}

	player, err := cfg.db.GetPlayerByUsername(r.Context(), newPlayer.Username)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Wrong Username or password", err)
		return
	}

	if err = auth.CompareHashPassword(player.HashedPassword, newPlayer.Password); err != nil {
		respondWithError(w, http.StatusBadRequest, "Wrong username or Password", err)
		return
	}

	refreshToken, err := cfg.db.GetRefreshTokenFromPlayerId(r.Context(), player.ID)
	if err != nil || time.Now().After(refreshToken.ExpiresAt) || refreshToken.RevokedAt.Valid  {
		cfg.db.DeleteRefreshToken(r.Context(), refreshToken.Token)

		refreshToken, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token: 		auth.MakeRefreshToken(),
			PlayerID: 	player.ID,
		})
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create new refresh token", err)
		return
	}

	token, err := auth.MakeJWT(player.ID, cfg.tokenSecret, time.Minute * 60)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create new JWT token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Player{
		Id: 		  player.ID,
		CreatedAt: 	  player.CreatedAt,
		Username: 	  player.Username,
		RefreshToken: refreshToken.Token,
		Token: 		  token,
	})
}

func (cfg *apiConfig) handlerGetPlayer(w http.ResponseWriter, r *http.Request) {
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

	player, err := cfg.db.GetPlayerByPlayerId(r.Context(), playerId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Player not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Player{
		Id: 		player.ID,
		Username: 	player.Username,
		Avatar: 	player.Avatar.String,
	})
}

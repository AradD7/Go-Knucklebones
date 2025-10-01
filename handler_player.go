package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/AradD7/Go-Knuclebones/internal/auth"
	"github.com/AradD7/Go-Knuclebones/internal/database"
	"github.com/AradD7/Go-Knuclebones/internal/verification"
	"github.com/google/uuid"
)

type Player struct {
	Id 			  uuid.UUID	`json:"id"`
	CreatedAt	  time.Time `json:"created_at"`
	Username 	  string  	`json:"username"`
	RefreshToken  string 	`json:"refresh_token"`
	Token 		  string 	`json:"token"`
	Avatar 		  string 	`json:"avatar"`
	DisplayName   string 	`json:"display_name"`
	Email 		  string 	`json:"email"`
	EmailVerified bool 		`json:"email_verified"`
}

func (cfg *apiConfig) handlerNewPlayer(w http.ResponseWriter, r *http.Request) {
	type createPlayerParams struct {
		Username  string `json:"username"`
		Email 	  string `json:"email"`
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
		HashedPassword: sql.NullString{
			Valid:  true,
			String: hashPassword,
		},
		Email: sql.NullString{
			Valid: true,
			String: newPlayer.Email,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to add the player to DB", err)
		return
	}

	token, hash := verification.GenerateVerificationToken()
	_, err = cfg.db.CreateVerificationToken(r.Context(), database.CreateVerificationTokenParams{
		TokenHash: 	hash,
		PlayerID: 	player.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate verification token", err)
		return
	}

	go verification.SendVerificationEmail(newPlayer.Email, token)

	respondWithJSON(w, http.StatusCreated, nil)
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

	if err = auth.CompareHashPassword(player.HashedPassword.String, newPlayer.Password); err != nil {
		respondWithError(w, http.StatusBadRequest, "Wrong username or Password", err)
		return
	}

	if !player.EmailVerified.Bool {
		respondWithError(w, http.StatusForbidden, "Player needs verification", err)
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
		Id: 		   player.ID,
		CreatedAt: 	   player.CreatedAt,
		Username: 	   player.Username,
		RefreshToken:  refreshToken.Token,
		Token: 		   token,
		Avatar: 	   player.Avatar.String,
		DisplayName:   player.DisplayName.String,
		EmailVerified: player.EmailVerified.Bool,
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
		Id: 		 player.ID,
		Username: 	 player.Username,
		Avatar: 	 player.Avatar.String,
		DisplayName: player.DisplayName.String,
	})
}

func (cfg *apiConfig) handlerUpdateProfile(w http.ResponseWriter, r *http.Request) {
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

	type paramaters struct {
		DisplayName string `json:"display_name"`
		Avatar 		string `json:"avatar"`
	}

	decoder := json.NewDecoder(r.Body)
	var params paramaters
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to decode the json data", err)
		return
	}

	if err = cfg.db.UpdateProfile(r.Context(), database.UpdateProfileParams{
		ID: 		 playerId,
		DisplayName: sql.NullString{
			Valid: 	params.DisplayName != "",
			String: params.DisplayName,
		},
		Avatar: 	 sql.NullString{
			Valid:  params.Avatar != "",
			String: params.Avatar,
		},
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update profile", err)
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

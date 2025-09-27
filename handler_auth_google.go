package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/AradD7/Go-Knuclebones/internal/auth"
	"github.com/AradD7/Go-Knuclebones/internal/database"
	"google.golang.org/api/idtoken"
)

func (cfg *apiConfig) handlerAuthGoogle(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		IdToken string `json:"id_token"`
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to read json data", err)
		return
	}

	payload, err := idtoken.Validate(r.Context(), params.IdToken, cfg.googleClientId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invlid id token", err)
		return
	}

	googleId  := payload.Subject
	email 	  := payload.Claims["email"].(string)
	firstName := payload.Claims["given_name"].(string)

	player, err := cfg.db.GetPlayerByGoogleId(r.Context(), sql.NullString{
		Valid:  true,
		String: googleId,
	})
	if err != nil {
		player, err = cfg.db.CreatePlayerWithGoogle(r.Context(), database.CreatePlayerWithGoogleParams{
			Username: strings.Split(email, "@")[0],
			GoogleID: sql.NullString{
				Valid:  true,
				String: googleId,
			},
			Email: sql.NullString{
				Valid:  true,
				String: email,
			},
			DisplayName: sql.NullString{
				Valid:  firstName != "",
				String: firstName,
			},
		})
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Failed to find account and to create one", err)
			return
		}
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
		Avatar: 	  player.Avatar.String,
		DisplayName:  player.DisplayName.String,
	})
}

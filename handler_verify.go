package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/AradD7/Go-Knuclebones/internal/auth"
	"github.com/AradD7/Go-Knuclebones/internal/database"
	"github.com/AradD7/Go-Knuclebones/internal/verification"
)

func (cfg *apiConfig) handlerVerifyEmail(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Token string `json:"token"`
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	// Hash the token to find it in DB
	hasher := sha256.New()
	hasher.Write([]byte(params.Token))
	tokenHash := hex.EncodeToString(hasher.Sum(nil))

	// Get verification token
	verification, err := cfg.db.GetVerificationToken(r.Context(), tokenHash)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid token", err)
		return
	}

	if !verification.ExpiresAt.After(time.Now().UTC()) {
		cfg.db.DeleteVerificationToken(r.Context(), tokenHash)
		respondWithError(w, http.StatusBadRequest, "Token expired", nil)
		return
	}

	// Mark email as verified
	err = cfg.db.VerifyPlayerEmail(r.Context(), verification.PlayerID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to verify email", err)
		return
	}

	// Delete used token
	cfg.db.DeleteVerificationToken(r.Context(), tokenHash)

	// Get player info
	player, err := cfg.db.GetPlayerByPlayerId(r.Context(), verification.PlayerID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get player", err)
		return
	}

	// Create tokens and log them in
	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:    auth.MakeRefreshToken(),
		PlayerID: player.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create new refreshToken token", err)
		return
	}

	accessToken, err := auth.MakeJWT(player.ID, cfg.tokenSecret, time.Minute*60)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create new JWT token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Player{
		Id:           player.ID,
		Email:        player.Email.String,
		DisplayName:  player.DisplayName.String,
		Avatar:       player.Avatar.String,
		RefreshToken: refreshToken.Token,
		Token:        accessToken,
	})
}

func (cfg *apiConfig) handlerResendVerification(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email     string `json:"email"`
		Useraname string `json:"username"`
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	player, err := cfg.db.GetPlayerByEmail(r.Context(), sql.NullString{
		Valid:  true,
		String: params.Email,
	})
	if err != nil {
		player, err = cfg.db.GetPlayerByUsername(r.Context(), params.Useraname)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "player not found", err)
			return
		}
	}

	if player.EmailVerified.Bool {
		respondWithError(w, http.StatusBadRequest, "Email already verified", nil)
		return
	}

	// Check for existing valid token
	existing, _ := cfg.db.GetVerificationTokenByPlayerId(r.Context(), player.ID)
	if existing.TokenHash != "" && existing.ExpiresAt.After(time.Now().UTC()) {
		// don't send if token less than half hour old
		if time.Now().UTC().Sub(existing.CreatedAt) < 30*time.Minute {
			respondWithError(w, http.StatusTooManyRequests, "Please wait before requesting another email", nil)
			return
		}
		// Delete old token
		cfg.db.DeleteVerificationToken(r.Context(), existing.TokenHash)
	}

	// Generate new token
	token, hash := verification.GenerateVerificationToken()
	cfg.db.CreateVerificationToken(r.Context(), database.CreateVerificationTokenParams{
		TokenHash: hash,
		PlayerID:  player.ID,
	})

	// Send email
	go verification.SendVerificationEmail(player.Email.String, token)

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Verification email sent",
	})
}

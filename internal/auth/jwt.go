package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "knucklebones-access"
)

func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: string(TokenTypeAccess),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: userId.String(),
	})

	return newToken.SignedString([]byte(tokenSecret))
}


func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	tokenClaim := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaim,
		func(jwt *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		})
	if err != nil {
		return uuid.Nil, err
	}

	userId, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	if issuer, err := token.Claims.GetIssuer(); err != nil {
		return uuid.Nil, fmt.Errorf("failed to get jwt issuer: %v", err)
	} else if issuer != string(TokenTypeAccess) {
		return uuid.Nil, fmt.Errorf("Failed to validate JWT")
	}
	return uuid.Parse(userId)
}


func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("Authorization not found in the header")
	}

	splitAuth := strings.Split(auth, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", fmt.Errorf("Bearer not found in the Authorization header")
	}
	return splitAuth[1], nil
}

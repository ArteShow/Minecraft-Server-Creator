package jwtutil

import (
	"time"

	"github.com/ArteShow/Minecraft-Server-Creator/services/auth-service/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID string, ttl time.Duration) (string, error) {
	cfg, err := config.Read()
	if err != nil {
		return "", err
	}

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

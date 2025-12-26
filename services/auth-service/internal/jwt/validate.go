package jwtutil

import (
	"github.com/ArteShow/Minecraft-Server-Creator/services/auth-service/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateToken(tokenString string) (*Claims, error) {
	cfg, err := config.Read()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			return cfg.JWTSecret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
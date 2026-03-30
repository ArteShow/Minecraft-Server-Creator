package jwt

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID string `json:"user_id"`
	AdminID string `json:"admin_id"`
	jwt.RegisteredClaims
}
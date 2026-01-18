package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ArteShow/Minecraft-Server-Creator/services/api-gateway/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const OwnerIDKey contextKey = "owner_id"

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cfg, err := config.Read()
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return 
			}

			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTKey), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			userID, ok := claims["user_id"].(string)
			if !ok || userID == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), OwnerIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

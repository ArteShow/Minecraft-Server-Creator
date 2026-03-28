package middleware

import (
	"net/http"
	"strings"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/gateway/internal/config"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg, err := config.Read()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return 
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return 
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == ""{
			w.WriteHeader(http.StatusUnauthorized)
			return 
		}

		r.Header.Set("X-User-ID", userID)
		r.Header.Set("Authorization", "Bearer "+parts[1])
		next.ServeHTTP(w, r)
	})
}
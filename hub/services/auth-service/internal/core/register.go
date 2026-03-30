package core

import (
	"fmt"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/auth-service/internal/config"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/auth-service/internal/repository"
)

func RegisterUser(username, password, email, jwt, Type string) (string, error) {
	cfg, err := config.Read()
	if err != nil {
		return "", err
	}

	if Type == "user" {
		userID, err := repository.CreateUser(username, password, email)
		if err != nil {
			return "", err
		}
		
		return userID, nil
	} else if Type == "admin" && jwt == cfg.JWTSecret {
		userID, err := repository.CreateAdmin(username, password)
		if err != nil {
			return "", err
		}

		return userID, nil
	}

	return "", fmt.Errorf("invalid user type or JWT")
}
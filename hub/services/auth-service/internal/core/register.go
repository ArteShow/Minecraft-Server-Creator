package core

import "github.com/ArteShow/Minecraft-Server-Creator/hub/services/auth-service/internal/repository"

func RegisterUser(username, password, email string) (string, error) {
	userID, err := repository.CreateUser(username, password, email)
	if err != nil {
		return "", err
	}

	return userID, nil
}
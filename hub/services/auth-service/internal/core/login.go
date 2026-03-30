package core

import (
	"errors"
	"time"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/auth-service/internal/jwt"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/auth-service/internal/repository"
)

const JWTTTL = time.Hour * 24

func LoginUser(username, password string) (string, error) {
	var token string
	var err error
	userID, ok := repository.CheckUserLogin(username, password)
	if !ok {
		adminID, ok := repository.CheckAdminLogin(username, password)
		if !ok {
			return "", errors.New("no user found with username: " + username + " and password: " + password)
		}

		token, err = jwt.GenerateAdminToken(adminID, JWTTTL)
		if err != nil {
			return "", err
		}
	} else {
		token, err = jwt.GenerateToken(userID, JWTTTL)
		if err != nil {
			return "", err
		}
	}

	return token, err
}

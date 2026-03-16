package repository

import (
	_ "database/sql"
	"errors"

	"github.com/ArteShow/Minecraft-Server-Creator/user-service/internal/database"
	"github.com/ArteShow/Minecraft-Server-Creator/user-service/pkg/id"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(username, password, email string) (string, error) {
	db, err := database.Connect()
	if err != nil {
		return "", err
	}
	defer db.Close()

	userID := id.GenerateID()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	
	_, err = db.Exec(
		"INSERT INTO users (id, email, password, username) VALUES ($1, $2, $3, $4)",
		userID, email, hash, username,
	)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func LoginUser(username, password string) (string, error) {
	db, err := database.Connect()
	if err != nil {
		return "", err
	}
	defer db.Close()

	row := db.QueryRow("SELECT password, id FROM users WHERE username = $1", username)

	var hashPassword, userID string
	err = row.Scan(&hashPassword, &userID)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		return "", errors.New("invalid username or password")
	}

	return userID, nil
}
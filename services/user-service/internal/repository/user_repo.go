package repository

import (
	"database/sql"

	"github.com/ArteShow/Minecraft-Server-Creator/user-service/internal/database"
	"github.com/ArteShow/Minecraft-Server-Creator/user-service/pkg/id"
)

func CreateUser(username, password, email string) (string, error) {
	db, err := database.Connect()
	if err != nil {
		return "", err
	}
	defer db.Close()

	userID := id.GenerateID()

	_, err = db.Exec(
		"INSERT INTO users (id, email, password, username) VALUES ($1, $2, $3, $4)",
		userID, email, password, username,
	)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func DeleteUser(userID string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		"DELETE FROM users WHERE id = $1",
		userID,
	)
	return err
}

func GetID(username string) (string, error) {
	db, err := database.Connect()
	if err != nil {
		return "", err
	}
	defer db.Close()

	var userID string

	err = db.QueryRow(
		"SELECT id FROM users WHERE username = $1",
		username,
	).Scan(&userID)

	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return userID, nil
}

func GetPassword(userID string) (string, error) {
	db, err := database.Connect()
	if err != nil {
		return "", err
	}
	defer db.Close()

	var password string

	err = db.QueryRow(
		"SELECT password FROM users WHERE id = $1",
		userID,
	).Scan(&password)

	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return password, nil
}

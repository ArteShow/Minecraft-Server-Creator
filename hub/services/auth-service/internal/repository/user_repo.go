package repository

import (
	"time"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/auth-service/internal/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(username, password, email string) (string, error) {
	db, err := database.Connect()
	if err != nil {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	id := uuid.New().String()

	_, err = db.Exec(
		`INSERT INTO users (id, username, password, email, created_at) 
		 VALUES ($1, $2, $3, $4, $5)`,
		id,
		username,
		string(hash),
		email,
		time.Now(),
	)

	return id, err
}

func CheckUserLogin(username, password string) (string, bool) {
	db, err := database.Connect()
	if err != nil {
		return "", false
	}

	var hash, id string

	err = db.QueryRow(
		`SELECT password, id FROM users WHERE username = $1`,
		username,
	).Scan(&hash, &id)

	if err != nil {
		return "", false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return id, err == nil
}

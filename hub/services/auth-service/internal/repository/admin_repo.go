package repository

import (
	_ "database/sql"
	"time"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/auth-service/internal/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateAdmin(username, password string) (string, error) {
	db, err := database.Connect()
	if err != nil {
		return "", err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	id := uuid.New().String()
	_, err = db.Exec(`INSERT INTO admins (id, username, password, created_at) VALUES ($1, $2, $3, $4)`,
		id,
		username,
		string(hash),
		time.Now(),
	)
	return id, err
}

func CheckAdminLogin(username, password string) (string, bool) {
	db, err := database.Connect()
	if err != nil {
		return "", false
	}

	var hash, id string
	err = db.QueryRow(`SELECT password, id FROM admins WHERE username = $1`, username).Scan(&hash, &id)
	if err != nil {
		return "", false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return id, err == nil
}

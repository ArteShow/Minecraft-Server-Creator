package repository

import (
	"encoding/json"
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

func AddBundle(userID string, bundleName string, value int64) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	query := `
		UPDATE users
		SET bundles = jsonb_set(
			COALESCE(bundles, '{}'::jsonb),
			ARRAY[$1],
			to_jsonb($2::bigint),
			true
		)
		WHERE id = $3
	`
	_, err = db.Exec(query, bundleName, value, userID)
	return err
}

func GetBundles(userID string) (map[string]int64, error) {
	db, err := database.Connect()
	if err != nil {
		return map[string]int64{}, err
	}

	var bundlesJSON []byte
	if err = db.QueryRow(`SELECT bundles FROM users WHERE id = $1`, userID).Scan(&bundlesJSON); err != nil {
		return nil, err
	}

	bundles := make(map[string]int64)
	if len(bundlesJSON) == 0 {
		return bundles, nil
	}

	err = json.Unmarshal(bundlesJSON, &bundles)
	return bundles, err
}
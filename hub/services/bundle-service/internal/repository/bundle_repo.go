package repository

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/bundle-service/internal/database"
)

func CreateBundleKey(userID, bundle string) (string, error) {
	db, err := database.Connect()
	if err != nil {
		return "", err
	}

	bytes := make([]byte, 16)
	_, err = rand.Read(bytes)
	if err != nil {
		return "", err
	}
	key := hex.EncodeToString(bytes)

	_, err = db.Exec(
		`INSERT INTO bundle_keys (key, user_id, bundle) VALUES ($1, $2, $3)`,
		key, userID, bundle,
	)

	return key, err
}

func UseBundleKey(key string) (bool, string, string, error) {
	db, err := database.Connect()
	if err != nil {
		return false, "", "", err
	}

	tx, err := db.Begin()
	if err != nil {
		return false, "", "", err
	}
	defer tx.Rollback()

	var userID, bundle string

	err = tx.QueryRow(
		`SELECT user_id, bundle FROM bundle_keys WHERE key = $1 FOR UPDATE`,
		key,
	).Scan(&userID, &bundle)

	if err != nil {
		return false, "", "", nil
	}

	_, err = tx.Exec(`DELETE FROM bundle_keys WHERE key = $1`, key)
	if err != nil {
		return false, "", "", err
	}

	if err = tx.Commit(); err != nil {
		return false, "", "", err
	}

	return true, userID, bundle, nil
}
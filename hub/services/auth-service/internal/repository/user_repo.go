package repository

import (
	_ "database/sql"
	"encoding/json"
	"errors"
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
		return "", err }
	id := uuid.New().String() 
	_, err = db.Exec(`INSERT INTO users (id, username, password, email, created_at) VALUES ($1, $2, $3, $4, $5)`,
		id,
		username,
		string(hash),
		email, time.Now(),
	) 
	return id, err
} 

func CheckUserLogin(username, password string) (string, bool) { 
	db, err := database.Connect()
	if err != nil { 
		return "", false 
	}

	var hash, id string 
	err = db.QueryRow( `SELECT password, id FROM users WHERE username = $1`, username, ).Scan(&hash, &id)
	if err != nil {
		return "", false 
	} 
	
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) 
	return id, err == nil }

func AddBundle(userID string, bundleName string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	var bundlesJSON []byte
	err = db.QueryRow(`SELECT bundles FROM users WHERE id = $1`, userID).Scan(&bundlesJSON)
	if err != nil {
		return err
	}

	bundles := make(map[string]int64)
	if len(bundlesJSON) > 0 {
		if err := json.Unmarshal(bundlesJSON, &bundles); err != nil {
			return err
		}
	}

	bundles[bundleName]++ 

	newJSON, err := json.Marshal(bundles)
	if err != nil {
		return err
	}

	_, err = db.Exec(`UPDATE users SET bundles = $1 WHERE id = $2`, newJSON, userID)
	return err
}

func DeleteBundle(userID string, bundleName string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	var bundlesJSON []byte
	err = db.QueryRow(`SELECT bundles FROM users WHERE id = $1`, userID).Scan(&bundlesJSON)
	if err != nil {
		return err
	}

	if len(bundlesJSON) == 0 {
		return errors.New("no bundles found")
	}

	bundles := make(map[string]int64)
	if err := json.Unmarshal(bundlesJSON, &bundles); err != nil {
		return err
	}

	if count, exists := bundles[bundleName]; exists {
		if count <= 1 {
			delete(bundles, bundleName) 
		} else {
			bundles[bundleName]-- 
		}
	} else {
		return errors.New("bundle not found")
	}

	newJSON, err := json.Marshal(bundles)
	if err != nil {
		return err
	}

	_, err = db.Exec(`UPDATE users SET bundles = $1 WHERE id = $2`, newJSON, userID)
	return err
}

func GetBundles(userID string) (map[string]int64, error) {
	db, err := database.Connect()
	if err != nil {
		return map[string]int64{}, err
	}

	var bundlesJSON []byte
	err = db.QueryRow(`SELECT bundles FROM users WHERE id = $1`, userID).Scan(&bundlesJSON)
	if err != nil {
		return nil, err
	}

	bundles := make(map[string]int64)
	if len(bundlesJSON) == 0 {
		return bundles, nil
	}

	err = json.Unmarshal(bundlesJSON, &bundles)
	return bundles, err
}
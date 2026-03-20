package repository

import (
	"database/sql"
	"time"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/network-service/internal/database"
	"github.com/google/uuid"
)

type HostMetadata struct {
	ID        string
	IP        string
	CreatedAt time.Time
}

func CreateServerFunction(ip string) (string, error) {
	serverID := uuid.NewString()

	db, err := database.Connect()
	if err != nil {
		return "", err
	}
	defer db.Close()

	_, err = db.Exec(
		`INSERT INTO network (id, ip)
		 VALUES ($1, $2)`,
		serverID, ip,
	)

	return serverID, err
}

func GetServerMetadataByID(serverID string) (HostMetadata, error) {
	db, err := database.Connect()
	if err != nil {
		return HostMetadata{}, err
	}
	defer db.Close()

	var metadata HostMetadata

	err = db.QueryRow(
		`SELECT id, ip, created_at
		 FROM network
		 WHERE id = $1`,
		serverID,
	).Scan(&metadata.ID, &metadata.IP, &metadata.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return HostMetadata{}, sql.ErrNoRows
		}
		return HostMetadata{}, err
	}

	return metadata, nil
}

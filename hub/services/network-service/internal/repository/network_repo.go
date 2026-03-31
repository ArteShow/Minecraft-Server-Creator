package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/network-service/internal/database"
	"github.com/google/uuid"
)

type HostMetadata struct {
	ID        string
	IP        string
	CreatedAt time.Time
}

func CreateServerFunction(ip, hostServerID string) (string, error) {
	if ip == "" {
		return "", fmt.Errorf("ip is required")
	}

	serverID := hostServerID
	if serverID == "" {
		serverID = uuid.NewString()
	}

	db, err := database.Connect()
	if err != nil {
		return "", err
	}
	defer db.Close()

	// Verify the host metadata row exists before inserting (network.id is a FK → hosts.id)
	var exists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM hosts WHERE id = $1)`, serverID).Scan(&exists)
	if err != nil {
		return "", fmt.Errorf("could not verify host: %w", err)
	}
	if !exists {
		return "", fmt.Errorf("host ID %q not found — create host metadata first, then use that ID here", serverID)
	}

	_, err = db.Exec(
		`INSERT INTO network (id, ip) VALUES ($1, $2)
		 ON CONFLICT (id) DO UPDATE SET ip = EXCLUDED.ip`,
		serverID, ip,
	)
	if err != nil {
		return "", fmt.Errorf("failed to register network host: %w", err)
	}

	return serverID, nil
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

package repository

import (
	"database/sql"

	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/internal/database"
)

func CreateServer(serverID, containerID, ownerID string, port int) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		`INSERT INTO servers (id, owner_id, container_id, port)
		 VALUES ($1, $2, $3, $4)`,
		serverID, ownerID, containerID, port,
	)

	return err
}

func GetUsersServersIDs(ownerID string) ([]string, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(
		`SELECT id FROM servers WHERE owner_id = $1`,
		ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var serverIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		serverIDs = append(serverIDs, id)
	}

	return serverIDs, rows.Err()
}

func GetUsersContainerIDs(ownerID string) ([]string, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(
		`SELECT container_id FROM servers WHERE owner_id = $1`,
		ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var containerIDs []string
	for rows.Next() {
		var id sql.NullString
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		if id.Valid {
			containerIDs = append(containerIDs, id.String)
		}
	}

	return containerIDs, rows.Err()
}

func DeleteServer(serverID string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		`DELETE FROM servers WHERE id = $1`,
		serverID,
	)

	return err
}

package repository

import (
	"database/sql"

	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/internal/database"
)

func CreateServer(serverID, ownerID string, port int) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		`INSERT INTO servers (id, owner_id, port)
		 VALUES ($1, $2, $3)`,
		serverID, ownerID, port,
	)

	return err
}

func AddContainerIDToServer(serverID, containerID string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		`UPDATE servers
		 SET container_id = $1
		 WHERE id = $2`,
		containerID, serverID,
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

func IsContainerOwnedByUser(containerID, ownerID string) (bool, error) {
	db, err := database.Connect()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow(
		`SELECT EXISTS (
			SELECT 1 FROM servers
			WHERE container_id = $1 AND owner_id = $2
		)`,
		containerID, ownerID,
	).Scan(&exists)

	return exists, err
}


func IsServerOwnedByUser(serverID, ownerID string) (bool, error) {
	db, err := database.Connect()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow(
		`SELECT EXISTS (
			SELECT 1 FROM servers
			WHERE id = $1 AND owner_id = $2
		)`,
		serverID, ownerID,
	).Scan(&exists)

	return exists, err
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

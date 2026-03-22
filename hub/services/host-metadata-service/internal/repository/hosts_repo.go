package repository

import (
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/host-metadata-service/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Host struct {
	ID        string   `json:"host_server_id"`
	Servers   []string `json:"server_ids"`
	CreatedAt string   `json:"created_at"`
}

func Get() ([]Host, error) {
	db, err := database.Connect()
	if err != nil {
		return []Host{}, err
	}

	var hosts []Host

	rows, err := db.Query(`SELECT id, servers, created_at FROM hosts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var host Host
		var servers pq.StringArray

		err := rows.Scan(&host.ID, &servers, &host.CreatedAt)
		if err != nil {
			return nil, err
		}
		host.Servers = []string(servers)
		hosts = append(hosts, host)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}

func Create(servers []string) (string, error) {
	db, err := database.Connect()
	if err != nil {
		return "", err
	}

	if servers == nil {
		servers = []string{}
	}

	id := uuid.NewString()

	_, err = db.Exec(
		`INSERT INTO hosts (id, servers) VALUES ($1, $2)`,
		id,
		pq.Array(servers),
	)
	return id, err
}

func Delete(id string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	_, err = db.Exec(`DELETE FROM hosts WHERE id = $1`, id)
	return err
}

func AddServer(hostID string, serverID string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`UPDATE hosts SET servers = array_append(servers, $1) WHERE id = $2`,
		serverID,
		hostID,
	)
	return err
}

func RemoveServer(hostID string, serverID string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`UPDATE hosts SET servers = array_remove(servers, $1) WHERE id = $2`,
		serverID,
		hostID,
	)
	return err
}

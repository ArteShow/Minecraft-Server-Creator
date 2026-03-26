package repository

import (
	"fmt"
	"strconv"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/host-metadata-service/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Host struct {
	ID        string   `json:"host_server_id"`
	Servers   []string `json:"server_ids"`
	RAM       string   `json:"RAM"`
	Cores     string   `json:"cpu_cores"`
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

func Create(servers []string, ram, cores string) (string, error) {
	db, err := database.Connect()
	if err != nil {
		return "", err
	}

	if servers == nil {
		servers = []string{}
	}

	id := uuid.NewString()

	_, err = db.Exec(
		`INSERT INTO hosts (id, servers, ram, cores) VALUES ($1, $2, $3, $4)`,
		id,
		pq.Array(servers),
		ram,
		cores,
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

func SubtractRAM(hostID string, ramToSubtract string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	var currentRAM string
	err = db.QueryRow(`SELECT ram FROM hosts WHERE id = $1`, hostID).Scan(&currentRAM)
	if err != nil {
		return err
	}

	current, err := strconv.ParseInt(currentRAM, 10, 64)
	if err != nil {
		return err
	}

	subtract, err := strconv.ParseInt(ramToSubtract, 10, 64)
	if err != nil {
		return err
	}

	newRAM := current - subtract

	_, err = db.Exec(
		`UPDATE hosts SET ram = $1 WHERE id = $2`,
		fmt.Sprintf("%d", newRAM),
		hostID,
	)
	return err
}

func SubtractCores(hostID string, coresToSubtract string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	var currentCores string
	err = db.QueryRow(`SELECT cores FROM hosts WHERE id = $1`, hostID).Scan(&currentCores)
	if err != nil {
		return err
	}

	current, err := strconv.ParseInt(currentCores, 10, 64)
	if err != nil {
		return err
	}

	subtract, err := strconv.ParseInt(coresToSubtract, 10, 64)
	if err != nil {
		return err
	}

	newCores := current - subtract

	_, err = db.Exec(
		`UPDATE hosts SET cores = $1 WHERE id = $2`,
		fmt.Sprintf("%d", newCores),
		hostID,
	)
	return err
}

func GetRAM(hostID string) (string, error) {
	db, err := database.Connect()
	if err != nil {
		return "", err
	}

	var ram string
	err = db.QueryRow(`SELECT ram FROM hosts WHERE id = $1`, hostID).Scan(&ram)
	if err != nil {
		return "", err
	}

	return ram, nil
}

func GetCores(hostID string) (string, error) {
	db, err := database.Connect()
	if err != nil {
		return "", err
	}

	var cores string
	err = db.QueryRow(`SELECT cores FROM hosts WHERE id = $1`, hostID).Scan(&cores)
	if err != nil {
		return "", err
	}

	return cores, nil
}

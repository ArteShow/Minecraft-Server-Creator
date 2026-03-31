package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/host-metadata-service/internal/database"
	"github.com/google/uuid"
)

type Host struct {
	ID        string         `json:"host_server_id"`
	Servers   map[string]int `json:"servers"`
	RAM       string         `json:"ram"`
	Cores     string         `json:"cpu_cores"`
	CreatedAt string         `json:"created_at"`
}

func loadServers(serversJSON []byte) (map[string]int, error) {
	servers := make(map[string]int)
	if len(serversJSON) == 0 {
		return servers, nil
	}

	if err := json.Unmarshal(serversJSON, &servers); err == nil {
		return servers, nil
	}

	legacyServers := make(map[string][]int)
	if err := json.Unmarshal(serversJSON, &legacyServers); err != nil {
		return nil, err
	}

	for serverID, ports := range legacyServers {
		if len(ports) == 0 {
			servers[serverID] = 0
			continue
		}
		servers[serverID] = ports[len(ports)-1]
	}

	return servers, nil
}

func saveServers(tx *sql.Tx, hostID string, servers map[string]int) error {
	newJSON, err := json.Marshal(servers)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE hosts SET servers = $1 WHERE id = $2`, newJSON, hostID)
	return err
}

func loadServersForUpdate(tx *sql.Tx, hostID string) (map[string]int, error) {
	var serversJSON []byte
	if err := tx.QueryRow(`SELECT servers FROM hosts WHERE id = $1 FOR UPDATE`, hostID).Scan(&serversJSON); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("host %s not found", hostID)
		}
		return nil, err
	}

	return loadServers(serversJSON)
}

func Get() ([]Host, error) {
	db, err := database.Connect()
	if err != nil {
		return []Host{}, err
	}

	var hosts []Host

	rows, err := db.Query(`SELECT id, servers, ram, cores, created_at FROM hosts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var host Host
		var serversJSON []byte
		var ram int
		var cores int

		err := rows.Scan(&host.ID, &serversJSON, &ram, &cores, &host.CreatedAt)
		if err != nil {
			return nil, err
		}

		host.Servers, err = loadServers(serversJSON)
		if err != nil {
			return nil, err
		}
		host.RAM = strconv.Itoa(ram)
		host.Cores = strconv.Itoa(cores)
		hosts = append(hosts, host)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}

func Create(servers map[string]int, ram, cores string) (string, error) {
	db, err := database.Connect()
	if err != nil {
		return "", err
	}

	id := uuid.NewString()

	serversJSON, err := json.Marshal(servers)
	if err != nil {
		return "", err
	}

	_, err = db.Exec(
		`INSERT INTO hosts (id, servers, ram, cores) VALUES ($1, $2, $3, $4)`,
		id,
		serversJSON,
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

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	servers, err := loadServersForUpdate(tx, hostID)
	if err != nil {
		return err
	}

	if _, exists := servers[serverID]; !exists {
		servers[serverID] = 0
	}

	if err = saveServers(tx, hostID, servers); err != nil {
		return err
	}

	return tx.Commit()
}

func RemoveServer(hostID string, serverID string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	servers, err := loadServersForUpdate(tx, hostID)
	if err != nil {
		return err
	}

	delete(servers, serverID)

	if err = saveServers(tx, hostID, servers); err != nil {
		return err
	}

	return tx.Commit()
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

	var ram int
	err = db.QueryRow(`SELECT ram FROM hosts WHERE id = $1`, hostID).Scan(&ram)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(ram), nil
}

func GetCores(hostID string) (string, error) {
	db, err := database.Connect()
	if err != nil {
		return "", err
	}

	var cores int
	err = db.QueryRow(`SELECT cores FROM hosts WHERE id = $1`, hostID).Scan(&cores)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(cores), nil
}

func AddPortToServer(hostID string, serverID string, port int) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	servers, err := loadServersForUpdate(tx, hostID)
	if err != nil {
		return err
	}

	servers[serverID] = port

	if err = saveServers(tx, hostID, servers); err != nil {
		return err
	}

	return tx.Commit()
}

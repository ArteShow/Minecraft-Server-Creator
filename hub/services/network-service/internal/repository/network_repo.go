package repository

import (
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/network-service/internal/database"
	"github.com/google/uuid"
)

func CreateServerFuntion(ip string) error {
	serverID := uuid.NewString()

	db, err := database.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		`INSERT INTO hosts (id, ip)
		 VALUES ($1, $2)`,
		serverID, ip,
	)

	return err
}

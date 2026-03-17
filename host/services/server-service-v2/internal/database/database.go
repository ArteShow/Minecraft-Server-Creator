package database

import (
	"database/sql"

	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/internal/config"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	cfg, err := config.Read()
	if err != nil {
		return nil, err
	}

	connStr :=
		"host=" + cfg.DBHost +
			" port=" + cfg.DBPort +
			" user=" + cfg.DBUser +
			" password=" + cfg.DBPassword +
			" dbname=" + cfg.DBName +
			" sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
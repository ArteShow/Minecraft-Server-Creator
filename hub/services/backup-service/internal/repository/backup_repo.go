package repository

import (
	"context"
	"time"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/backup-service/internal/database"
	"github.com/google/uuid"
)

type Backup struct {
	BackupID  string
	ServerID  string
	UserID    string
	CreatedAt time.Time
}

func CreateBackup(ctx context.Context, b *Backup) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.ExecContext(ctx,
		`INSERT INTO backups (backup_id, server_id, user_id) VALUES ($1,$2,$3)`,
		uuid.NewString(), b.ServerID, b.UserID,
	)
	return err
}

func GetBackups(ctx context.Context, serverID string) ([]Backup, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx,
		`SELECT backup_id, server_id, user_id, created_at 
		 FROM backups 
		 WHERE server_id = $1
		 ORDER BY created_at DESC`,
		serverID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var backups []Backup
	for rows.Next() {
		var b Backup
		if err := rows.Scan(&b.BackupID, &b.ServerID, &b.UserID, &b.CreatedAt); err != nil {
			return nil, err
		}
		backups = append(backups, b)
	}

	return backups, nil
}

func DeleteBackup(ctx context.Context, backupID string) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.ExecContext(ctx, `DELETE FROM backups WHERE backup_id=$1`, backupID)
	return err
}

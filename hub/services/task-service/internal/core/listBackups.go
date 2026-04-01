package core

import (
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/client"
	backup "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/backup-service"
)

type BackupItem struct {
	BackupID string
	ServerID string
	UserID   string
}

func ListBackups(serverID, userID string) ([]BackupItem, error) {
	backupClient, err := client.NewBackupClient()
	if err != nil {
		return []BackupItem{}, err
	}
	defer backupClient.Close()

	resp, err := backupClient.GetBackup(&backup.GetBackupRequest{ServerID: serverID})
	if err != nil {
		return []BackupItem{}, err
	}

	items := make([]BackupItem, 0, len(resp.GetBackups()))
	for _, entry := range resp.GetBackups() {
		if userID != "" && entry.GetUserID() != userID {
			continue
		}

		items = append(items, BackupItem{
			BackupID: entry.GetBackup(),
			ServerID: entry.GetServerID(),
			UserID:   entry.GetUserID(),
		})
	}

	return items, nil
}

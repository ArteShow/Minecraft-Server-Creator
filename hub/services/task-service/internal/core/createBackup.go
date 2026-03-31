package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/client"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/config"
	backup "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/backup-service"
	host "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/host-metadata-service"
	network "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/network-service"
)

func CreateBackup(serverID, token, userID, bundle string) error {
	bundles, err := config.GetBundles()
	if err != nil {
		return err
	}

	cfg, err := config.Read()
	if err != nil {
		return err
	}

	backupClient, err := client.NewBackupClient()
	if err != nil {
		return err
	}

	backups, err := backupClient.GetBackup(&backup.GetBackupRequest{ServerID: serverID})
	if err != nil {
		return err
	}

	var backupCounter int
	for _, backup := range backups.GetBackups() {
		if backup.GetUserID() == userID {
			backupCounter++
		}
	}
	if backupCounter > bundles.Bundles[bundle].Backups {
		return fmt.Errorf("You are not allowed to have this amount of backups in this bundle %s", bundle)
	}

	hostClient, err := client.NewHostClient()
	if err != nil {
		return err
	}

	hosts, err := hostClient.GetAllHostServers(&host.GetAllHostServersRequest{})
	if err != nil {
		return err
	}

	hostID := SelecthostIDByServerID(serverID, hosts)
	if hostID == "" {
		return fmt.Errorf("server %s is not mapped to a host", serverID)
	}

	networkClient, err := client.NewNetworkClient()
	if err != nil {
		return err
	}

	ip, err := networkClient.GetServerMetadata(&network.GetServerMetadataRequest{ServerId: hostID})
	if err != nil {
		return err
	}

	requestBody := map[string]string{"server_id": serverID}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		"http://"+ip.Ip+":"+cfg.DefaultHostServerPort+"/server-service/backup/create",
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	_, err = backupClient.CreateBackup(&backup.CreateBackupRequest{ServerID: serverID, UserID: userID})
	if err != nil {
		return err
	}

	return nil
}

package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/client"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/config"
	backup "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/backup-service"
	host "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/host-metadata-service"
	network "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/network-service"
)

func DeleteBackup(serverID, token, backupID string) error {
	cfg, err := config.Read()
	if err != nil {
		return err
	}

	hostClient, err := client.NewHostClient()
	if err != nil {
		return err
	}
	defer hostClient.Close()

	backupClient, err := client.NewBackupClient()
	if err != nil {
		return err
	}
	defer backupClient.Close()

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
	defer networkClient.Close()

	ip, err := networkClient.GetServerMetadata(&network.GetServerMetadataRequest{ServerId: hostID})
	if err != nil {
		return err
	}

	targetIP := normalizeHostIP(ip.Ip)
	if targetIP == "" {
		return fmt.Errorf("host metadata returned empty IP for host %s", hostID)
	}

	requestBody := map[string]string{"server_id": serverID, "backup_id": backupID}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		"http://"+targetIP+":"+cfg.DefaultHostServerPort+"/server-service/backup/delete",
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if len(body) > 0 {
			return fmt.Errorf("host backup delete failed, status %d: %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("host backup delete failed, status %d", resp.StatusCode)
	}

	_, err = backupClient.DeleteBackup(&backup.DeleteBackupRequest{BackupID: backupID})
	if err != nil {
		return err
	}

	return nil
}

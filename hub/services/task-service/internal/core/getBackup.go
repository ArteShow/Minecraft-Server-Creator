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

func GetBackup(serverID, backupID, token string) ([]byte, error) {
	cfg, err := config.Read()
	if err != nil {
		return []byte{}, err
	}

	hostClient, err := client.NewHostClient()
	if err != nil {
		return []byte{}, err
	}
	defer hostClient.Close()

	hosts, err := hostClient.GetAllHostServers(&host.GetAllHostServersRequest{})
	if err != nil {
		return []byte{}, err
	}

	hostID := SelecthostIDByServerID(serverID, hosts)
	if hostID == "" {
		return []byte{}, fmt.Errorf("server %s is not mapped to a host", serverID)
	}

	networkClient, err := client.NewNetworkClient()
	if err != nil {
		return []byte{}, err
	}
	defer networkClient.Close()

	ip, err := networkClient.GetServerMetadata(&network.GetServerMetadataRequest{ServerId: hostID})
	if err != nil {
		return []byte{}, err
	}

	targetIP := normalizeHostIP(ip.Ip)
	if targetIP == "" {
		return []byte{}, fmt.Errorf("host metadata returned empty IP for host %s", hostID)
	}

	if backupID == "" {
		backupClient, clientErr := client.NewBackupClient()
		if clientErr != nil {
			return []byte{}, clientErr
		}
		defer backupClient.Close()

		resp, clientErr := backupClient.GetBackup(&backup.GetBackupRequest{ServerID: serverID})
		if clientErr != nil {
			return []byte{}, clientErr
		}
		if len(resp.GetBackups()) == 0 {
			return []byte{}, fmt.Errorf("no backups found for server %s", serverID)
		}
		backupID = resp.GetBackups()[0].GetBackup()
	}

	requestBody := map[string]string{"server_id": serverID, "backup_id": backupID}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return []byte{}, err
	}

	req, err := http.NewRequest(
		"POST",
		"http://"+targetIP+":"+cfg.DefaultHostServerPort+"/server-service/backup/get",
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if len(body) > 0 {
			return []byte{}, fmt.Errorf("host backup download failed, status %d: %s", resp.StatusCode, string(body))
		}
		return []byte{}, fmt.Errorf("host backup download failed, status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

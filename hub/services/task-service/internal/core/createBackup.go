package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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

	selectedBundle := ""
	bundleLimit := -1
	for name, cfg := range bundles.Bundles {
		if strings.EqualFold(name, strings.TrimSpace(bundle)) {
			selectedBundle = name
			bundleLimit = cfg.Backups
			break
		}
	}
	if selectedBundle == "" {
		return fmt.Errorf("unknown bundle %q", bundle)
	}

	if backupCounter >= bundleLimit {
		return fmt.Errorf("backup limit reached for bundle %s (%d)", selectedBundle, bundleLimit)
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

	targetIP := normalizeHostIP(ip.Ip)
	if targetIP == "" {
		return fmt.Errorf("host metadata returned empty IP for host %s", hostID)
	}

	knownBackupIDs := map[string]struct{}{}
	for _, item := range backups.GetBackups() {
		if item.GetUserID() == userID {
			knownBackupIDs[item.GetBackup()] = struct{}{}
		}
	}

	_, err = backupClient.CreateBackup(&backup.CreateBackupRequest{ServerID: serverID, UserID: userID})
	if err != nil {
		return err
	}

	createdList, err := backupClient.GetBackup(&backup.GetBackupRequest{ServerID: serverID})
	if err != nil {
		return err
	}

	backupID := ""
	for _, item := range createdList.GetBackups() {
		if item.GetUserID() != userID {
			continue
		}
		if _, exists := knownBackupIDs[item.GetBackup()]; exists {
			continue
		}
		backupID = item.GetBackup()
		break
	}

	if backupID == "" {
		return fmt.Errorf("failed to resolve created backup id for server %s", serverID)
	}

	requestBody := map[string]string{"server_id": serverID, "backup_id": backupID}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		_, _ = backupClient.DeleteBackup(&backup.DeleteBackupRequest{BackupID: backupID})
		return err
	}

	req, err := http.NewRequest(
		"POST",
		"http://"+targetIP+":"+cfg.DefaultHostServerPort+"/server-service/backup/create",
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		_, _ = backupClient.DeleteBackup(&backup.DeleteBackupRequest{BackupID: backupID})
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		_, _ = backupClient.DeleteBackup(&backup.DeleteBackupRequest{BackupID: backupID})
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		_, _ = backupClient.DeleteBackup(&backup.DeleteBackupRequest{BackupID: backupID})
		if len(body) > 0 {
			return fmt.Errorf("host backup create failed, status %d: %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("host backup create failed, status %d", resp.StatusCode)
	}

	return nil
}

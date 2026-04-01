package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/client"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/config"
	host "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/host-metadata-service"
	network "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/network-service"
)

func DeleteServer(serverID, token, ownerID string) error {
	cfg, err := config.Read()
	if err != nil {
		return err
	}

	hostClient, err := client.NewHostClient()
	if err != nil {
		return err
	}
	defer hostClient.Close()

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

	requestBody := map[string]string{"server_id": serverID}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		"http://"+targetIP+":"+cfg.DefaultHostServerPort+"/server-service/delete",
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Owner-ID", ownerID)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	_, err = hostClient.RemoveServerFromHost(&host.RemoveServerFromHostRequest{ServerId: serverID, HostServerId: hostID})
	if err != nil {
		return err
	}

	return nil
}

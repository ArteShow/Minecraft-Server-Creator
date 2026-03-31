package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/client"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/config"
	host "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/host-metadata-service"
	network "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/network-service"
)

func GetBackup(serverID, token string) ([]byte, error) {
	cfg, err := config.Read()
	if err != nil {
		return []byte{}, err
	}

	hostClient, err := client.NewHostClient()
	if err != nil {
		return []byte{}, err
	}

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

	ip, err := networkClient.GetServerMetadata(&network.GetServerMetadataRequest{ServerId: hostID})
	if err != nil {
		return []byte{}, err
	}

	requestBody := map[string]string{"server_id": serverID}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return []byte{}, err
	}

	req, err := http.NewRequest(
		"POST",
		"http://"+ip.Ip+":"+cfg.DefaultHostServerPort+"/server-service/backup/get",
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	return body, nil
}

package core

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/client"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/config"
	host "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/host-metadata-service"
	network "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/network-service"
)

func GetServerStats(key, serverID, token string) (string, error) {
	cfg, err := config.Read()
	if err != nil {
		return "", err
	}

	hostClient, err := client.NewHostClient()
	if err != nil {
		return "", err
	}

	hosts, err := hostClient.GetAllHostServers(&host.GetAllHostServersRequest{})
	if err != nil {
		return "", err
	}

	hostID := SelectHostWithFewestServers(*hosts)

	networkClient, err := client.NewNetworkClient()
	if err != nil {
		return "", err
	}

	serverMetadata, err := networkClient.GetServerMetadata(&network.GetServerMetadataRequest{ServerId: hostID})
	if err != nil {
		return "", err
	}

	requestBody := map[string]string{"key": key, "server_id": serverID}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		"http://"+serverMetadata.Ip+":"+cfg.DefaultHostServerPort+"/server/getServerStats",
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var respData GetServerStatsResponse
	if err = json.Unmarshal(body, &respData); err != nil {
		return "", err
	}

	return respData.Value, nil
}

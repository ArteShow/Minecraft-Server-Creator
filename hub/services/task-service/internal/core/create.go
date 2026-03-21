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

func CreateServer(version string) (string, int, error) {
	cfg, err := config.Read()
	if err != nil {
		return "", 0, err
	}

	hostClient, err := client.NewHostClient()
	if err != nil {
		return "", 0, err
	}

	servers, err := hostClient.GetAllHostServers(&host.GetAllHostServersRequest{})
	if err != nil {
		return "", 0, err
	}

	var minServers int = -1
	var selectedHostId string
	for _, host := range servers.Hosts {
		serverCount := len(host.Servers)
		if minServers == -1 || serverCount < minServers {
			minServers = serverCount
			selectedHostId = host.Id
		}
	}

	networkClient, err := client.NewNetworkClient()
	if err != nil {
		return "", 0, err
	}

	serverMetadata, err := networkClient.GetServerMetadata(&network.GetServerMetadataRequest{ServerId: selectedHostId})
	if err != nil {
		return "", 0, err
	}

	requestBody := map[string]string{"version": version}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", 0, err
	}

	resp, err := http.Post(
		"http://"+serverMetadata.Ip+":"+cfg.DefaultHostServerPort+"/server/create",
		"application/json",
		bytes.NewReader(jsonBody),
	)
	if err != nil || resp.StatusCode != http.StatusCreated {
		return "", 0, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	var respData CreateServerResponse
	if err = json.Unmarshal(body, &respData); err != nil {
		return "", 0, err
	}

	_, err = hostClient.AddServerToHost(&host.AddServerToHostRequest{HostServerId: selectedHostId, ServerId: respData.ServerID})
	if err != nil {
		return "", 0, err
	}

	return respData.ServerID, respData.Port, nil
}

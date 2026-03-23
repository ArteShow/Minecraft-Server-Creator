package core

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/client"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/config"
	host "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/host-metadata-service"
	network "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/network-service"
)

func DeleteServer(serverID, token string) error {
	cfg, err := config.Read()
	if err != nil {
		return err
	}

	hostClient, err := client.NewHostClient()
	if err != nil {
		return err
	}

	hosts, err := hostClient.GetAllHostServers(&host.GetAllHostServersRequest{})
	if err != nil {
		return err
	}

	var hostID string
	for _, host := range hosts.Hosts {
		for _, mcServerID := range host.Servers {
			if mcServerID == serverID {
				hostID = mcServerID
			}
		}
	}

	networkClient, err := client.NewNetworkClient()
	if err != nil {
		return err
	}

	ip, err := networkClient.GetServerMetadata(&network.GetServerMetadataRequest{ServerId: hostID})
	if err != nil {
		return err
	}

	requestBody := map[string]string{"server_id": hostID}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		"http://"+ip.Ip+":"+cfg.DefaultHostServerPort+"/server/delete",
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

	_, err = hostClient.RemoveServerFromHost(&host.RemoveServerFromHostRequest{ServerId: serverID, HostServerId: hostID})
	if err != nil {
		return err
	}

	return nil
}

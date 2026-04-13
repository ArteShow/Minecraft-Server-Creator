package core

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/client"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/config"
	host "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/host-metadata-service"
	network "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/network-service"
)

func InstallPlugin(serverID, token string, filename string, data []byte) error {
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

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("server_id", serverID)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return err
	}
	if _, err = io.Copy(part, bytes.NewReader(data)); err != nil {
		return err
	}
	writer.Close()

	req, err := http.NewRequest(
		"POST",
		"http://"+targetIP+":"+cfg.DefaultHostServerPort+"/server-service/plugin/install",
		body,
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		if len(respBody) > 0 {
			return fmt.Errorf("host plugin install failed, status %d: %s", resp.StatusCode, string(respBody))
		}
		return fmt.Errorf("host plugin install failed, status %d", resp.StatusCode)
	}

	return nil
}

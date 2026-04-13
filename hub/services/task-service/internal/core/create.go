package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/client"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/config"
	clientBundle "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/bundle-service"
	host "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/host-metadata-service"
	network "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/network-service"
)

const firstServerPort = 25565

func SelectHostWithFewestServers(servers *host.GetAllHostServersResponse) string {
	minServers := -1
	selectedHostID := ""
	for _, currentHost := range servers.GetHosts() {
		count := len(currentHost.Servers)
		if minServers == -1 || count < minServers {
			minServers = count
			selectedHostID = currentHost.Id
		}
	}
	return selectedHostID
}

func SelecthostIDByServerID(serverID string, servers *host.GetAllHostServersResponse) string {
	for _, currentHost := range servers.GetHosts() {
		for key := range currentHost.GetServers() {
			if key == serverID {
				return currentHost.GetId()
			}
		}
	}
	return ""
}

func GetNextPort(servers *host.GetAllHostServersResponse) int {
	highestPort := firstServerPort - 1
	for _, currentHost := range servers.GetHosts() {
		for _, value := range currentHost.GetServers() {
			if highestPort < int(value) {
				highestPort = int(value)
			}
		}
	}
	return highestPort + 1
}

func CreateServer(version, serverType, token, bundleKey, userID string) (string, int, error) {
	cfg, err := config.Read()
	if err != nil {
		return "", 0, err
	}

	bundles, err := config.GetBundles()
	if err != nil {
		return "", 0, err
	}

	hostClient, err := client.NewHostClient()
	if err != nil {
		return "", 0, err
	}
	defer hostClient.Close()

	networkClient, err := client.NewNetworkClient()
	if err != nil {
		return "", 0, err
	}
	defer networkClient.Close()

	bundleClient, err := client.NewBundleClient()
	if err != nil {
		return "", 0, err
	}
	defer bundleClient.Close()

	serversResp, err := hostClient.GetAllHostServers(&host.GetAllHostServersRequest{})
	if err != nil {
		return "", 0, err
	}

	hostID := SelectHostWithFewestServers(serversResp)
	if hostID == "" {
		return "", 0, fmt.Errorf("no hosts available")
	}

	ramResp, err := hostClient.GetRAM(&host.GetRAMRequest{HostServerId: hostID})
	if err != nil {
		return "", 0, err
	}
	availableRAM, err := strconv.Atoi(ramResp.GetRam())
	if err != nil {
		return "", 0, err
	}

	coresResp, err := hostClient.GetCores(&host.GetCoresRequest{HostServerId: hostID})
	if err != nil {
		return "", 0, err
	}
	availableCores, err := strconv.Atoi(coresResp.GetCores())
	if err != nil {
		return "", 0, err
	}

	bundle, err := bundleClient.DisableBundleKey(&clientBundle.DisableBundleKeyRequest{Key: bundleKey})
	if err != nil {
		return "", 0, err
	}

	_, err = bundleClient.RemoveBundle(&clientBundle.RemoveBundleRequest{UserID: userID, Bundle: bundle.GetBundle()})
	if err != nil {
		return "", 0, err
	}

	bundleData, ok := bundles.Bundles[bundle.GetBundle()]
	if !ok {
		return "", 0, fmt.Errorf("unknown bundle: %s", bundle.GetBundle())
	}

	if bundleData.RAM > availableRAM || bundleData.Cores > availableCores {
		return "", 0, fmt.Errorf("host %s does not have enough resources", hostID)
	}

	metaResp, err := networkClient.GetServerMetadata(&network.GetServerMetadataRequest{ServerId: hostID})
	if err != nil {
		return "", 0, err
	}

	targetIP := normalizeHostIP(metaResp.Ip)
	if targetIP == "" {
		return "", 0, fmt.Errorf("host metadata returned empty IP for host %s", hostID)
	}

	nextPort := GetNextPort(serversResp)
	bodyBytes, err := json.Marshal(map[string]interface{}{"version": version, "port": nextPort, "server_type": serverType})
	if err != nil {
		return "", 0, err
	}

	url := fmt.Sprintf("http://%s:%s/server-service/create", targetIP, cfg.DefaultHostServerPort)
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Owner-ID", userID)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		errBody, _ := io.ReadAll(resp.Body)
		if len(errBody) > 0 {
			return "", 0, fmt.Errorf("failed to create server, status: %d, response: %s", resp.StatusCode, string(errBody))
		}
		return "", 0, fmt.Errorf("failed to create server, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	var respData CreateServerResponse
	if err := json.Unmarshal(body, &respData); err != nil {
		return "", 0, err
	}
	if respData.Port == 0 {
		respData.Port = nextPort
	}

	_, err = hostClient.AddPortToServer(&host.AddPortToServerRequest{
		HostServerId: hostID,
		ServerId:     respData.ServerID,
		Port:         int32(respData.Port),
	})
	if err != nil {
		return "", 0, err
	}

	return respData.ServerID, respData.Port, nil
}

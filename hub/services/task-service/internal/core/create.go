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

func SelectHostWithFewestServers(servers host.GetAllHostServersResponse) string {
	minServers := -1
	selectedHostId := ""
	for _, h := range servers.GetHosts() {
		count := len(h.Servers)
		if minServers == -1 || count < minServers {
			minServers = count
			selectedHostId = h.Id
		}
	}
	return selectedHostId
}

func CreateServer(version, token, bundle_key, userID string) (string, int, error) {
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

	networkClient, err := client.NewNetworkClient()
	if err != nil {
		return "", 0, err
	}

	bundleClient, err := client.NewBundleClient()
	if err != nil {
		return "", 0, err
	}

	serversResp, err := hostClient.GetAllHostServers(&host.GetAllHostServersRequest{})
	if err != nil {
		return "", 0, err
	}

	hostID := SelectHostWithFewestServers(*serversResp)
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

	bundle, err := bundleClient.DisableBundleKey(&clientBundle.DisableBundleKeyRequest{Key: bundle_key})
	if err != nil {
		return "", 0, err
	}

	_, err = bundleClient.RemoveBundle(&clientBundle.RemoveBundleRequest{UserID: userID, Bundle: bundle.GetBundle()})

	bundleData, ok := bundles.Bundles[bundle.GetBundle()]
	if !ok {
		return "", 0, fmt.Errorf("unknown bundle: %s", bundle)
	}

	if bundleData.RAM > availableRAM || bundleData.Cores > availableCores {
		return "", 0, fmt.Errorf("host %s does not have enough resources", hostID)
	}

	metaResp, err := networkClient.GetServerMetadata(&network.GetServerMetadataRequest{ServerId: hostID})
	if err != nil {
		return "", 0, err
	}

	bodyBytes, err := json.Marshal(map[string]string{"version": version})
	if err != nil {
		return "", 0, err
	}

	url := fmt.Sprintf("http://%s:%s/server/create", metaResp.Ip, cfg.DefaultHostServerPort)
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
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

	_, err = hostClient.AddServerToHost(&host.AddServerToHostRequest{
		HostServerId: hostID,
		ServerId:     respData.ServerID,
	})
	if err != nil {
		return "", 0, err
	}

	return respData.ServerID, respData.Port, nil
}
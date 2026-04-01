package core

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/client"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/config"
	host "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/host-metadata-service"
	network "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/network-service"
)

type ConsoleSnapshot struct {
	Console       string `json:"console"`
	OnlinePlayers int    `json:"online_players"`
}

var (
	joinedPattern = regexp.MustCompile(`(?i)\]:\s([^\s]+) joined the game`)
	leftPattern   = regexp.MustCompile(`(?i)\]:\s([^\s]+) left the game`)
)

func GetServerConsoleSnapshot(serverID, token string, tail int) (*ConsoleSnapshot, error) {
	cfg, err := config.Read()
	if err != nil {
		return nil, err
	}

	hostClient, err := client.NewHostClient()
	if err != nil {
		return nil, err
	}

	hosts, err := hostClient.GetAllHostServers(&host.GetAllHostServersRequest{})
	if err != nil {
		return nil, err
	}

	hostID := SelecthostIDByServerID(serverID, hosts)
	if hostID == "" {
		return nil, fmt.Errorf("server %s is not mapped to a host", serverID)
	}

	networkClient, err := client.NewNetworkClient()
	if err != nil {
		return nil, err
	}

	serverMetadata, err := networkClient.GetServerMetadata(&network.GetServerMetadataRequest{ServerId: hostID})
	if err != nil {
		return nil, err
	}

	targetIP := normalizeHostIP(serverMetadata.Ip)
	if targetIP == "" {
		return nil, fmt.Errorf("host metadata returned empty IP for host %s", hostID)
	}

	url := fmt.Sprintf("http://%s:%s/server-service/console?server_id=%s&tail=%d", targetIP, cfg.DefaultHostServerPort, serverID, tail)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		msg := strings.TrimSpace(string(body))
		if msg == "" {
			msg = "failed to fetch console"
		}
		return nil, fmt.Errorf("console request failed: %s", msg)
	}

	console := string(body)
	return &ConsoleSnapshot{
		Console:       console,
		OnlinePlayers: countOnlinePlayers(console),
	}, nil
}

func countOnlinePlayers(console string) int {
	online := map[string]struct{}{}
	lines := strings.Split(console, "\n")

	for _, line := range lines {
		if match := joinedPattern.FindStringSubmatch(line); len(match) == 2 {
			online[match[1]] = struct{}{}
		}
		if match := leftPattern.FindStringSubmatch(line); len(match) == 2 {
			delete(online, match[1])
		}
	}

	return len(online)
}

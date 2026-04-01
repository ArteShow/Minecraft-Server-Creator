package core

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
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
	loginPattern  = regexp.MustCompile(`(?i)\]:\s([^\s\[]+)\[.*logged in with entity id`)
	leftPattern   = regexp.MustCompile(`(?i)\]:\s([^\s]+) left the game`)
	listPattern   = regexp.MustCompile(`(?i)there are\s+(\d+)\s+of a max`)
)

func ResolveServerTarget(serverID string) (string, error) {
	hostClient, err := client.NewHostClient()
	if err != nil {
		return "", err
	}
	defer hostClient.Close()

	hosts, err := hostClient.GetAllHostServers(&host.GetAllHostServersRequest{})
	if err != nil {
		return "", err
	}

	hostID := SelecthostIDByServerID(serverID, hosts)
	if hostID == "" {
		return "", fmt.Errorf("server %s is not mapped to a host", serverID)
	}

	networkClient, err := client.NewNetworkClient()
	if err != nil {
		return "", err
	}
	defer networkClient.Close()

	serverMetadata, err := networkClient.GetServerMetadata(&network.GetServerMetadataRequest{ServerId: hostID})
	if err != nil {
		return "", err
	}

	targetIP := normalizeHostIP(serverMetadata.Ip)
	if targetIP == "" {
		return "", fmt.Errorf("host metadata returned empty IP for host %s", hostID)
	}

	return targetIP, nil
}

func GetServerConsoleSnapshotFromTarget(targetIP, serverID, token string, tail int) (*ConsoleSnapshot, error) {
	cfg, err := config.Read()
	if err != nil {
		return nil, err
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

func GetServerConsoleSnapshot(serverID, token string, tail int) (*ConsoleSnapshot, error) {
	targetIP, err := ResolveServerTarget(serverID)
	if err != nil {
		return nil, err
	}

	return GetServerConsoleSnapshotFromTarget(targetIP, serverID, token, tail)
}

func countOnlinePlayers(console string) int {
	online := map[string]struct{}{}
	lines := strings.Split(console, "\n")
	lastKnown := -1

	for _, line := range lines {
		if match := listPattern.FindStringSubmatch(line); len(match) == 2 {
			if parsed, err := strconv.Atoi(match[1]); err == nil {
				lastKnown = parsed
			}
		}
		if match := joinedPattern.FindStringSubmatch(line); len(match) == 2 {
			online[match[1]] = struct{}{}
		}
		if match := loginPattern.FindStringSubmatch(line); len(match) == 2 {
			online[match[1]] = struct{}{}
		}
		if match := leftPattern.FindStringSubmatch(line); len(match) == 2 {
			delete(online, match[1])
		}
	}

	if len(online) > 0 {
		return len(online)
	}
	if lastKnown >= 0 {
		return lastKnown
	}
	return 0
}

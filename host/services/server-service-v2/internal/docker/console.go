package docker

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
)

func (ds *DockerService) GetConsoleLogs(serverID string, tail int) (string, error) {
	if strings.TrimSpace(serverID) == "" {
		return "", fmt.Errorf("serverID is required")
	}
	if tail <= 0 {
		tail = 200
	}

	ctx := context.Background()
	reader, err := ds.client.ContainerLogs(ctx, "mc_container_"+serverID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       strconv.Itoa(tail),
	})
	if err != nil {
		return "", err
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	clean := strings.ReplaceAll(string(content), "\x00", "")
	return strings.TrimSpace(clean), nil
}

func (ds *DockerService) SendConsoleCommand(serverID, command string) error {
	if strings.TrimSpace(serverID) == "" {
		return fmt.Errorf("serverID is required")
	}
	cmd := strings.TrimSpace(command)
	if cmd == "" {
		return fmt.Errorf("command is required")
	}

	ctx := context.Background()
	attached, err := ds.client.ContainerAttach(ctx, "mc_container_"+serverID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Logs:   false,
	})
	if err != nil {
		return err
	}
	defer attached.Close()

	if _, err = io.WriteString(attached.Conn, cmd+"\n"); err != nil {
		return err
	}

	if c, ok := attached.Conn.(interface{ CloseWrite() error }); ok {
		_ = c.CloseWrite()
	}

	return nil
}

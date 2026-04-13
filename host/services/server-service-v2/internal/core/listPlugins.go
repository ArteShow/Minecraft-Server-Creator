package core

import (
	"fmt"
	"strings"
)

func (s *Server) ListPlugins(serverID string) ([]string, error) {
	if serverID == "" {
		return nil, fmt.Errorf("server_id is required")
	}

	output, err := s.DockerService.ExecuteCommandInVolume(
		"mc_"+serverID,
		[]string{"sh", "-c", "if [ -d /data/plugins ]; then ls -1 /data/plugins; fi"},
	)
	if err != nil {
		return nil, fmt.Errorf("list plugins: %w", err)
	}

	lines := strings.Split(strings.ReplaceAll(output, "\r\n", "\n"), "\n")
	plugins := make([]string, 0, len(lines))
	for _, line := range lines {
		name := strings.TrimSpace(line)
		if name == "" || name == ".keep" {
			continue
		}
		if strings.HasSuffix(strings.ToLower(name), ".jar") {
			plugins = append(plugins, name)
		}
	}

	return plugins, nil
}

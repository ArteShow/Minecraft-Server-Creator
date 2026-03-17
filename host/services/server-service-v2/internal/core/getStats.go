package core

import "github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/internal/stats"

func (s *Server) GetStats(key, serverID string) (string, error) {
	file, err := s.DockerService.GetFileFromVolume(serverID, "/data", "stats.json")
	if err != nil {
		return "", err
	}

	value, err := stats.GetValue(key, file)
	if err != nil {
		return "", err
	}

	return value, nil
}
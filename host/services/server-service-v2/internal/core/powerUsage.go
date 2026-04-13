package core

import (
	"errors"

	"github.com/ArteShow/Minecraft-Server-Creator/host/services/server-service-v2/internal/repository"
)

type PowerUsage struct {
	Online     bool    `json:"online"`
	CPUPercent float64 `json:"cpu_percent"`
	RAMUsedMB  float64 `json:"ram_used_mb"`
	RAMLimitMB float64 `json:"ram_limit_mb"`
	RAMPercent float64 `json:"ram_percent"`
}

func (s *Server) GetPowerUsage(serverID, ownerID string) (PowerUsage, error) {
	ok, err := repository.IsServerOwnedByUser(serverID, ownerID)
	if err != nil {
		return PowerUsage{}, errors.New("failed to verify server ownership: " + err.Error())
	}
	if !ok {
		return PowerUsage{}, errors.New("user with id: " + ownerID + " is not the owner of this server")
	}

	containerID, err := repository.GetServerContainerID(serverID)
	if err != nil {
		return PowerUsage{}, err
	}
	if containerID == "" {
		return PowerUsage{Online: false}, nil
	}

	cpuPercent, usedBytes, limitBytes, err := s.DockerService.GetContainerResourceUsage(containerID)
	if err != nil {
		return PowerUsage{}, err
	}

	ramUsedMB := float64(usedBytes) / (1024.0 * 1024.0)
	ramLimitMB := float64(limitBytes) / (1024.0 * 1024.0)
	ramPercent := 0.0
	if ramLimitMB > 0 {
		ramPercent = (ramUsedMB / ramLimitMB) * 100.0
	}

	return PowerUsage{
		Online:     true,
		CPUPercent: cpuPercent,
		RAMUsedMB:  ramUsedMB,
		RAMLimitMB: ramLimitMB,
		RAMPercent: ramPercent,
	}, nil
}

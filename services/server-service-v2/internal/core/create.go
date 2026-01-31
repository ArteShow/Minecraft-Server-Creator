package core

import (
	"fmt"

	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/internal/repository"
	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/pkg/eula"
	get_version "github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/pkg/version"
	"github.com/google/uuid"
)

func (s *Server) CreateServer(version, ownerID string) (string, int, error) {
	id := uuid.NewString()

	if err := s.DockerService.CreateVolume(id); err != nil {
		return "", 0, err
	}

	jar, err := get_version.GetServerJar(version)
	if err != nil {
		return "", 0, fmt.Errorf("failed to download jar: %w", err)
	}

	if err := s.DockerService.UploadToVolume(
		id,
		"/data",
		"server.jar",
		jar,
	); err != nil {
		return "", 0, err
	}

	eula, err := eula.Accept()
	if err != nil {
		return "", 0, err
	}

	if err = s.DockerService.UploadToVolume(
		id,
		"/data",
		"eula.txt",
		eula,
	); err != nil {
		return "", 0, err
	}

	port, err := repository.GetHighestPort()
	if err != nil {
		return "", 0, err
	}

	if port == 0 {
		port = 25564
	}

	if err = repository.CreateServer(id, ownerID, port+1); err != nil {
		return "", 0, err
	}

	return id, port+1, nil
}

package core

import (
	"fmt"

	"github.com/ArteShow/Minecraft-Server-Creator/host/services/server-service-v2/internal/repository"
	"github.com/ArteShow/Minecraft-Server-Creator/host/services/server-service-v2/internal/stats"
	"github.com/ArteShow/Minecraft-Server-Creator/host/services/server-service-v2/pkg/eula"
	get_version "github.com/ArteShow/Minecraft-Server-Creator/host/services/server-service-v2/pkg/version"
	"github.com/google/uuid"
)

func (s *Server) CreateServer(version, serverType, ownerID string, port int) (string, error) {
	id := uuid.NewString()

	if err := s.DockerService.CreateVolume(id); err != nil {
		return "", err
	}

	var jar []byte
	var err error
	switch serverType {
	case "Paper":
		jar, err = get_version.GetPaperJar(version)
	case "Spigot":
		jar, err = get_version.GetSpigotJar(version)
	default:
		jar, err = get_version.GetServerJar(version)
	}
	if err != nil {
		return "", fmt.Errorf("failed to download jar: %w", err)
	}

	if err := s.DockerService.UploadToVolume(
		id,
		"/data",
		"server.jar",
		jar,
	); err != nil {
		return "", err
	}

	eulaTxt, err := eula.Accept()
	if err != nil {
		return "", err
	}

	if err = s.DockerService.UploadToVolume(
		id,
		"/data",
		"eula.txt",
		eulaTxt,
	); err != nil {
		return "", err
	}

	templete, err := stats.CreateTemplete()
	if err != nil {
		return "", err
	}
	if err = s.DockerService.UploadToVolume(
		id,
		"/data",
		"stats.json",
		templete,
	); err != nil {
		return "", err
	}

	if err = repository.CreateServer(id, ownerID, port); err != nil {
		return "", err
	}

	return id, nil
}

package core

import (
	"fmt"

	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/pkg/eula"
	get_version "github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/pkg/version"
	"github.com/google/uuid"
)

func (s *Server )CreateServer(version string) (string, error){
	id := uuid.NewString()

	if err := s.DockerService.CreateVolume(id); err != nil {
		return "", err
	}

	jar, err := get_version.GetServerJar(version)
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

	eula, err := eula.Accept()
	if err != nil {
		return "", err
	}

	if err = s.DockerService.UploadToVolume(
		id,
		"/data",
		"eula.txt",
		eula,
	); err != nil {
		return "", err
	}

	return id, nil
}
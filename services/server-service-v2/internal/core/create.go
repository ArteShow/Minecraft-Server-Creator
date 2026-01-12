package core

import (
	"fmt"

	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/internal/docker"
	get_version "github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/pkg/version"
	"github.com/google/uuid"
)

func CreateServer(version string, ds *docker.DockerService) (string, error){
	id := uuid.NewString()

	if err := ds.CreateVolume(id); err != nil {
		return "", err
	}

	jar, err := get_version.GetServerJar(version)
	if err != nil {
		return "", fmt.Errorf("failed to download jar: %w", err)
	}

	if err := ds.UploadToVolume(
		id,
		"/data",
		"server.jar",
		jar,
	); err != nil {
		return "", err
	}

	return id, nil
}
package core

import (
	"errors"

	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/internal/repository"
)

func (s *Server) StartServer(serverID, ownerID string) (string, error) {
	ok, err := repository.IsServerOwnedByUser(serverID, ownerID)
	if err != nil || !ok {
		return "", errors.New("user with id: " + ownerID + " is not the owner of this server: +" + err.Error())
	}

	port, err := repository.GetServersPort(serverID)
	if err != nil {
		return "", err
	}

	conID, err := s.DockerService.StartServerContainer(serverID, "eclipse-temurin:21-jre-jammy", port, 25565)
	if err != nil {
		return "", err
	}

	s.Processes.Add(serverID, conID)
	if err = repository.AddContainerIDToServer(serverID, conID); err != nil {
		return "", err
	}

	return conID, nil
}

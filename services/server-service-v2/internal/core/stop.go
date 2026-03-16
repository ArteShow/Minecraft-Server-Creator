package core

import (
	"errors"

	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/internal/repository"
)

func (s *Server) StopServer(serverID, ownerID string) error {
	containerID, ok := s.Processes.Get(serverID)
	if !ok {
		return errors.New("failed to find the container id")
	}

	ok, err := repository.IsServerOwnedByUser(serverID, ownerID)
	if err != nil || !ok {
		return errors.New("user with id: " + ownerID + " is not the owner of this server: +" + err.Error())
	}
	
	ok, err = repository.IsContainerOwnedByUser(containerID, ownerID)
	if err != nil || !ok {
		return errors.New("user with id: " + ownerID + " is not the owner of this container: +" + err.Error())
	}

	if err = repository.RemoveContainerID(containerID, ownerID); err != nil {
		return err
	}
	s.Processes.Remove(serverID)

	if err := s.DockerService.StopContainer(containerID); err != nil {
		return err
	}

	return s.DockerService.RemoveContainer(containerID)
}

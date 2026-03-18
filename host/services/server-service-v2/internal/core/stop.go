package core

import (
	"errors"

	"github.com/ArteShow/Minecraft-Server-Creator/host/services/server-service-v2/internal/repository"
	"github.com/ArteShow/Minecraft-Server-Creator/host/services/server-service-v2/internal/stats"
)

func (s *Server) StopServer(serverID, ownerID string) error {
	containerID, ok := s.Processes.Get(serverID)
	if !ok {
		return errors.New("failed to find the container id")
	}

	ok, err := repository.IsServerOwnedByUser(serverID, ownerID)
	if err != nil {
		return errors.New("failed to verify server ownership: " + err.Error())
	}
	if !ok {
		return errors.New("user with id: " + ownerID + " is not the owner of this server")
	}

	ok, err = repository.IsContainerOwnedByUser(containerID, ownerID)
	if err != nil {
		return errors.New("failed to verify container ownership: " + err.Error())
	}
	if !ok {
		return errors.New("user with id: " + ownerID + " is not the owner of this container")
	}

	if err := s.DockerService.StopContainer(containerID); err != nil {
		return err
	}

	if err := s.DockerService.RemoveContainer(containerID); err != nil {
		return err
	}

	if err = repository.RemoveContainerID(containerID, ownerID); err != nil {
		return err
	}
	s.Processes.Remove(serverID)

	file, err := s.DockerService.GetFileFromVolume(serverID, "/data", "stats.json")
	if err != nil {
		return err
	}

	changes, err := stats.SetValue("Online", false, file)
	if err != nil {
		return err
	}

	if err = s.DockerService.DeleteFileFromVolume(serverID, "/data", "stats.json"); err != nil {
		return err
	}

	if err = s.DockerService.UploadToVolume(serverID, "/data", "stats.json", changes); err != nil {
		return err
	}

	return nil
}

package core

import (
	"errors"

	"github.com/ArteShow/Minecraft-Server-Creator/host/services/server-service-v2/internal/repository"
	"github.com/ArteShow/Minecraft-Server-Creator/host/services/server-service-v2/internal/stats"
)

func (s *Server) StartServer(serverID, ownerID, RAM string, cores int) error {
	ok, err := repository.IsServerOwnedByUser(serverID, ownerID)
	if err != nil {
		return errors.New("failed to verify server ownership: " + err.Error())
	}
	if !ok {
		return errors.New("user with id: " + ownerID + " is not the owner of this server")
	}

	port, err := repository.GetServersPort(serverID)
	if err != nil {
		return err
	}

	conID, err := s.DockerService.StartServerContainer(serverID, "eclipse-temurin:21-jre-jammy", RAM, cores, port, 25565)
	if err != nil {
		return err
	}

	s.Processes.Add(serverID, conID)
	if err = repository.AddContainerIDToServer(serverID, conID); err != nil {
		return err
	}

	file, err := s.DockerService.GetFileFromVolume(serverID, "/data", "stats.json")
	if err != nil {
		return err
	}

	changes, err := stats.SetValue("Online", true, file)
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

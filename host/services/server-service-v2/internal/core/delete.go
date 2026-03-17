package core

import (
	"errors"

	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/internal/repository"
)

func (s *Server) DeleteServer(serverID, ownerID string) error {
	ok, err := repository.IsServerOwnedByUser(serverID, ownerID)
	if err != nil || !ok {
		return errors.New("user with id: " + ownerID + " is not the owner of this server: +" + err.Error())
	}

	if err = repository.DeleteServer(serverID); err != nil {
		return err
	}
	
	return s.DockerService.DeleteVolume(serverID)
}
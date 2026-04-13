package core

import "fmt"

func (s *Server) UploadWorld(world []byte, serverID string) error {
	// Delete existing world folder to make room for the new one
	if err := s.DockerService.DeleteFolderFromVolume(serverID, "/data/world"); err != nil {
		// Ignore error if folder doesn't exist yet
		if err.Error() != "" {
			_ = err // Log but don't fail
		}
	}

	// Upload and extract the world folder archive
	if err := s.DockerService.UploadFolderToVolume(serverID, "/data", world); err != nil {
		return fmt.Errorf("failed to upload world: %w", err)
	}

	return nil
}
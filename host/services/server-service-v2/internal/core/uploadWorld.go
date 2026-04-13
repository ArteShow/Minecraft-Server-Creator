package core

func (s *Server) UploadWorld(world []byte, serverID string) error {
	if err := s.DockerService.DeleteFolderFromVolume(serverID, "data/world"); err != nil {
		return err
	}

	if err := s.DockerService.UploadFolderToVolume(serverID, "data/world", world); err != nil {
		return err
	}

	return nil
}
package core

func (s *Server) UploadBackup(backup []byte, filename, serverID string) error {
	if err := s.DockerService.DeleteVolume("mc_" + serverID); err != nil {
		return err
	}

	if err := s.DockerService.CreateVolume("mc_" + serverID); err != nil {
		return err
	}

	if err := s.DockerService.UploadToVolume("mc_"+serverID, "data/", filename, backup); err != nil {
		return err
	}

	_, err := s.DockerService.ExecuteCommandInVolume("mc_"+serverID, []string{"sh", "-c", "tar -xvf /data/filename.tar -C /data"})
	if err != nil {
		return err
	}

	return nil
}
package core

func (s *Server) StopServer(serverID, containerID string) error {
	s.Processes.Remove(serverID)

	if err := s.DockerService.StopContainer(containerID); err != nil {
		return err
	}

	return s.DockerService.RemoveContainer(containerID)
}

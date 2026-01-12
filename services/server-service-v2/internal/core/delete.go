package core

func (s *Server) DeleteServer(serverID string) error {
	return s.DockerService.DeleteVolume(serverID)
}
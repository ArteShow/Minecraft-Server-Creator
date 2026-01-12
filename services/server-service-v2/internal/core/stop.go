package core

func (s *Server) StopServer(serverID, containerID string) error {
	s.Processes.Remove(serverID)
	return s.DockerService.ExecInContainer(
		containerID,
		[]string{
			"sh",
			"-c",
			"echo stop > /proc/1/fd/0",
		},
	)
}

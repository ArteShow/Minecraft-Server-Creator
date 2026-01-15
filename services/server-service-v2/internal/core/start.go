package core

func (s *Server) StartServer(serverID string) (string, error) {
	conID, err := s.DockerService.StartServerContainer(serverID, "eclipse-temurin:21-jre-jammy", 25565)
	if err != nil {
		return "", err
	}

	s.Processes.Add(serverID, conID)
	return conID, nil
}
package core

func (s *Service) StartServer(id string, port int) (string, error) {
	return s.ds.StartServerContainer(id, "eclipse-temurin:21-jre", port)
}
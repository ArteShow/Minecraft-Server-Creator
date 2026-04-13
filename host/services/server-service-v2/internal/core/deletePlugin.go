package core

import "fmt"

func (s *Server) DeletePlugin(serverID, fileName string) error {
	if serverID == "" {
		return fmt.Errorf("server_id is required")
	}
	if fileName == "" {
		return fmt.Errorf("file_name is required")
	}

	if err := s.DockerService.DeleteFileFromVolume(serverID, "/data/plugins", fileName); err != nil {
		return fmt.Errorf("delete plugin: %w", err)
	}
	return nil
}

package core

import "fmt"

// InstallPlugin writes a plugin jar into the server's /data/plugins/ directory.
func (s *Server) InstallPlugin(serverID, filename string, data []byte) error {
	if serverID == "" {
		return fmt.Errorf("server_id is required")
	}
	if filename == "" {
		return fmt.Errorf("filename is required")
	}
	if len(data) == 0 {
		return fmt.Errorf("plugin file is empty")
	}

	if err := s.DockerService.UploadToVolume(serverID, "/data/plugins", filename, data); err != nil {
		return fmt.Errorf("upload plugin to volume: %w", err)
	}
	return nil
}

package core

import "fmt"

func (s *Server) UploadBackup(backup []byte, filename, serverID string) error {
	if filename == "" {
		filename = "backup.tar.gz"
	}

	volumeName := "mc_" + serverID

	// Upload the backup archive to the volume
	if err := s.DockerService.UploadToVolume(serverID, "/data", filename, backup); err != nil {
		return fmt.Errorf("failed to upload backup file: %w", err)
	}

	// Extract the backup archive to /data (replacing existing world data)
	cmd := fmt.Sprintf("cd /data && tar -xzf %s && rm %s", filename, filename)
	_, err := s.DockerService.ExecuteCommandInVolume(volumeName, []string{"sh", "-c", cmd})
	if err != nil {
		return fmt.Errorf("failed to extract backup: %w", err)
	}

	return nil
}
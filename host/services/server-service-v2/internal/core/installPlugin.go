package core

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

const maxSupportedPluginClassMajor = 65 // Java 21

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

	if err := validatePluginJavaCompatibility(data); err != nil {
		return err
	}

	if err := s.DockerService.UploadToVolume(serverID, "/data/plugins", filename, data); err != nil {
		return fmt.Errorf("upload plugin to volume: %w", err)
	}
	return nil
}

func validatePluginJavaCompatibility(data []byte) error {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("invalid plugin jar: %w", err)
	}

	mainClassPath := findPluginMainClassPath(reader.File)
	if mainClassPath == "" {
		return nil
	}

	classFile := findFile(reader.File, mainClassPath)
	if classFile == nil {
		return nil
	}

	major, err := readClassMajorVersion(classFile)
	if err != nil {
		return fmt.Errorf("read plugin class version: %w", err)
	}

	if major > maxSupportedPluginClassMajor {
		return fmt.Errorf(
			"plugin requires Java %d (class file version %d), but server runtime supports up to Java %d",
			javaVersionFromClassMajor(major),
			major,
			javaVersionFromClassMajor(maxSupportedPluginClassMajor),
		)
	}

	return nil
}

func findPluginMainClassPath(files []*zip.File) string {
	for _, name := range []string{"plugin.yml", "paper-plugin.yml"} {
		descriptor := findFile(files, name)
		if descriptor == nil {
			continue
		}

		openFile, err := descriptor.Open()
		if err != nil {
			continue
		}
		descriptorData, readErr := io.ReadAll(openFile)
		_ = openFile.Close()
		if readErr != nil {
			continue
		}

		for _, rawLine := range strings.Split(string(descriptorData), "\n") {
			line := strings.TrimSpace(rawLine)
			if !strings.HasPrefix(line, "main:") {
				continue
			}
			mainClass := strings.TrimSpace(strings.TrimPrefix(line, "main:"))
			mainClass = strings.Trim(mainClass, `"'`)
			if mainClass == "" {
				continue
			}
			return strings.ReplaceAll(mainClass, ".", "/") + ".class"
		}
	}

	return ""
}

func findFile(files []*zip.File, name string) *zip.File {
	for _, file := range files {
		if file.Name == name {
			return file
		}
	}
	return nil
}

func readClassMajorVersion(file *zip.File) (uint16, error) {
	openFile, err := file.Open()
	if err != nil {
		return 0, err
	}
	defer openFile.Close()

	header := make([]byte, 8)
	if _, err = io.ReadFull(openFile, header); err != nil {
		return 0, err
	}

	if binary.BigEndian.Uint32(header[0:4]) != 0xCAFEBABE {
		return 0, fmt.Errorf("invalid class file magic")
	}

	return binary.BigEndian.Uint16(header[6:8]), nil
}

func javaVersionFromClassMajor(major uint16) int {
	if major <= 44 {
		return 0
	}
	return int(major) - 44
}

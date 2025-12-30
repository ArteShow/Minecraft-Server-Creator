package server

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service/internal/config"
	eulaacceptor "github.com/ArteShow/Minecraft-Server-Creator/services/server-service/pkg/eula_acceptor"
	getjar "github.com/ArteShow/Minecraft-Server-Creator/services/server-service/pkg/get_jar"
	idgenerator "github.com/ArteShow/Minecraft-Server-Creator/services/server-service/pkg/id_generator"
)

func CreateServer(version string) (string, error) {
	id := idgenerator.GenerateServerID()
	err := os.MkdirAll("./servers/"+id, 0755)
	if err != nil {
		return "", err
	}

	err = getjar.GetServerJar(version, "./servers/"+id+"/")
	if err != nil {
		return "", err
	}

	err = eulaacceptor.WriteEULA("./servers/" + id + "/")
	if err != nil {
		return "", err
	}

	return id, nil
}

func StartServer(serverID string) (*exec.Cmd, error) {
	serverPath := filepath.Join("servers", serverID)

	if _, err := os.Stat(serverPath); err != nil {
		return nil, fmt.Errorf("server folder not found")
	}

	if _, err := os.Stat(serverPath + "/server.jar"); err != nil {
		return nil, fmt.Errorf("server.jar not found")
	}

	cfg, err := config.Read()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(
		"java",
		"-Xms"+cfg.StartRAM,
		"-Xmx"+cfg.RunRAM,
		"-jar",
		"server.jar",
		"nogui",
	)

	cmd.Dir = serverPath + "/"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return cmd, nil
}

func StopServer(serverID string, cmd *exec.Cmd) error {
	if cmd == nil {
		return fmt.Errorf("cmd is nil")
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	_, err = stdin.Write([]byte("stop\n"))
	if err != nil {
		return err
	}

	cmd.Wait()
	return nil
}

func DeleteServer(serverID string) error {
	serverPath := filepath.Join("servers", serverID)
	err := os.RemoveAll(serverPath)
	if err != nil {
		return err
	}

	return nil
}

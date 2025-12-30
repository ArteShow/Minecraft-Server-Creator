package server

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service/internal/config"
	eulaacceptor "github.com/ArteShow/Minecraft-Server-Creator/services/server-service/pkg/eula_acceptor"
	getjar "github.com/ArteShow/Minecraft-Server-Creator/services/server-service/pkg/get_jar"
	idgenerator "github.com/ArteShow/Minecraft-Server-Creator/services/server-service/pkg/id_generator"
)

type ServerStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type ServerProcess struct {
	Cmd   *exec.Cmd
	Stdin io.WriteCloser
}

func CreateServer(version string) (string, error) {
	cwd, _ := os.Getwd()
	fmt.Println("Current working directory:", cwd)

	id := idgenerator.GenerateServerID()
	serverPath := filepath.Join("servers", id)

	fmt.Println("Creating folder:", serverPath)
	if err := os.MkdirAll(serverPath, 0755); err != nil {
		return "", err
	}

	statusFile := filepath.Join(serverPath, "status.json")
	status := ServerStatus{Status: "downloading"}
	saveStatus(statusFile, status)

	go func() {
		err := getjar.GetServerJar(version, serverPath)
		if err != nil {
			status = ServerStatus{Status: "error", Error: err.Error()}
		} else {
			_ = eulaacceptor.WriteEULA(serverPath)
			status = ServerStatus{Status: "ready"}
		}
		saveStatus(statusFile, status)
	}()

	return id, nil
}

func StartServer(serverID string) (*ServerProcess, error) {
	serverPath := filepath.Join("servers", serverID)
	statusFile := filepath.Join(serverPath, "status.json")

	var status ServerStatus
	for {
		status, _ = loadStatus(statusFile)
		if status.Status == "ready" {
			break
		}
		if status.Status == "error" {
			return nil, fmt.Errorf("server error: %s", status.Error)
		}
		time.Sleep(500 * time.Millisecond)
	}

	jarPath := filepath.Join(serverPath, "server.jar")
	if _, err := os.Stat(jarPath); err != nil {
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
	cmd.Dir = serverPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &ServerProcess{Cmd: cmd, Stdin: stdin}, nil
}

func StopServer(p *ServerProcess) error {
	if p == nil || p.Cmd == nil || p.Stdin == nil {
		return fmt.Errorf("server process not running")
	}

	_, err := p.Stdin.Write([]byte("stop\n"))
	if err != nil {
		return err
	}

	return p.Cmd.Wait()
}

func DeleteServer(serverID string) error {
	serverPath := filepath.Join("servers", serverID)
	return os.RemoveAll(serverPath)
}

func saveStatus(path string, status ServerStatus) {
	data, _ := json.Marshal(status)
	_ = os.WriteFile(path, data, 0644)
}

func loadStatus(path string) (ServerStatus, error) {
	var status ServerStatus
	data, err := os.ReadFile(path)
	if err != nil {
		return status, err
	}
	err = json.Unmarshal(data, &status)
	return status, err
}

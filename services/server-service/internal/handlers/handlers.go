package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service/internal/config"
	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service/internal/models"
	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service/internal/server"
	eulaacceptor "github.com/ArteShow/Minecraft-Server-Creator/services/server-service/pkg/eula_acceptor"
)

func CreateServer(w http.ResponseWriter, r *http.Request) {
	var req models.CreateServerRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := server.CreateServer(req.Version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := models.CreateServerResponse{
		ServerID: id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func StartServer(w http.ResponseWriter, r *http.Request) {
	var req models.StartServerRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	serverPath := "./servers/" + req.ServerID

	if _, err := os.Stat(serverPath); err != nil {
		http.Error(w, "server not found", http.StatusNotFound)
		return
	}

	cfg, err := config.Read()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := eulaacceptor.WriteEULA(serverPath + "/"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := models.StartServerResponse{
		Status: "running",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

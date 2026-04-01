package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type ConsoleCommandRequest struct {
	ServerID string `json:"server_id"`
	Command  string `json:"command"`
}

func (h *Handler) SendServerConsoleCommand(w http.ResponseWriter, r *http.Request) {
	var req ConsoleCommandRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.ServerID = strings.TrimSpace(req.ServerID)
	req.Command = strings.TrimSpace(req.Command)
	if req.ServerID == "" || req.Command == "" {
		http.Error(w, "server_id and command are required", http.StatusBadRequest)
		return
	}

	if err = h.Server.DockerService.SendConsoleCommand(req.ServerID, req.Command); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

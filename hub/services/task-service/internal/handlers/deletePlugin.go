package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/core"
)

type DeletePluginRequest struct {
	ServerID string `json:"server_id"`
	FileName string `json:"file_name"`
}

func DeletePlugin(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	token := parts[1]

	var req DeletePluginRequest
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

	if req.ServerID == "" || req.FileName == "" {
		http.Error(w, "server_id and file_name are required", http.StatusBadRequest)
		return
	}

	if err = core.DeletePlugin(req.ServerID, token, req.FileName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/core"
)

func GetPowerUsage(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimSpace(r.Header.Get("X-User-ID"))
	if userID == "" {
		http.Error(w, "userID header missing", http.StatusBadRequest)
		return
	}

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

	serverID := strings.TrimSpace(r.URL.Query().Get("server_id"))
	if serverID == "" {
		http.Error(w, "server_id is required", http.StatusBadRequest)
		return
	}

	usage, err := core.GetPowerUsage(serverID, token, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(usage); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

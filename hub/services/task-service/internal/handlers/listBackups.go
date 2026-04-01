package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/core"
)

func ListBackups(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "userID header missing", http.StatusBadRequest)
		return
	}

	var req GetBackupRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	backups, err := core.ListBackups(req.ServerID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items := make([]BackupItem, 0, len(backups))
	for _, entry := range backups {
		items = append(items, BackupItem{
			BackupID: entry.BackupID,
			ServerID: entry.ServerID,
			UserID:   entry.UserID,
		})
	}

	resp := ListBackupsResponse{Backups: items}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

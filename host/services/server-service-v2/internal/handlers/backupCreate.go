package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

func (h *Handler) CreateBackup(w http.ResponseWriter, r *http.Request) {
	var req CreateBackupRequest
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

	backupName := "mc_backup_" + req.ServerID
	if req.BackupID != "" {
		backupName += "_" + req.BackupID
	}

	if err = h.Server.DockerService.CreateBackup(req.ServerID, backupName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

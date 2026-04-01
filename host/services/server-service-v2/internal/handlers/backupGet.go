package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

func (h *Handler) Getbackup(w http.ResponseWriter, r *http.Request) {
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

	backupName := "mc_backup_" + req.ServerID
	if req.BackupID != "" {
		backupName += "_" + req.BackupID
	}

	backup, err := h.Server.DockerService.DownloadBackup(req.ServerID, backupName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\"backup.tar.gz\"")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(backup); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

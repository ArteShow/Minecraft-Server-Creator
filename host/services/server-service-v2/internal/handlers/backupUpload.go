package handlers

import (
	"io"
	"net/http"
)

func (h *Handler) UploadBackup(w http.ResponseWriter, r *http.Request) {
	serverID := r.FormValue("server_id")
	backupName := r.FormValue("backup_name")
	if serverID == "" {
		http.Error(w, "server_id is required", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "failed to read file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(data) == 0 {
		http.Error(w, "file is empty", http.StatusBadRequest)
		return
	}

	if err = h.Server.UploadBackup(data, backupName, serverID); err != nil {
		http.Error(w, "backup upload failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

func (h *Handler) DeleteBackup(w http.ResponseWriter, r * http.Request) {
	var req DeleteBackupRequest
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

	if err = h.Server.DockerService.DeleteBackup(req.ServerID, "mc_backup_"+req.ServerID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
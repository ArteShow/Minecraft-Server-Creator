package handlers

import (
	"io"
	"net/http"
)

func (h *Handler) InstallPlugin(w http.ResponseWriter, r *http.Request) {
	serverID := r.FormValue("server_id")
	if serverID == "" {
		http.Error(w, "server_id is required", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
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

	if err := h.Server.InstallPlugin(serverID, header.Filename, data); err != nil {
		http.Error(w, "plugin install failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

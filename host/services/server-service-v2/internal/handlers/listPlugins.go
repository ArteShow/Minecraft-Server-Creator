package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) ListPlugins(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("server_id")
	if serverID == "" {
		http.Error(w, "server_id is required", http.StatusBadRequest)
		return
	}

	plugins, err := h.Server.ListPlugins(serverID)
	if err != nil {
		http.Error(w, "plugin list failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"plugins": plugins})
}

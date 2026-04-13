package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (h *Handler) GetPowerUsage(w http.ResponseWriter, r *http.Request) {
	ownerID := strings.TrimSpace(r.Header.Get("X-Owner-ID"))
	if ownerID == "" {
		http.Error(w, "owner id missing", http.StatusBadRequest)
		return
	}

	serverID := strings.TrimSpace(r.URL.Query().Get("server_id"))
	if serverID == "" {
		http.Error(w, "server_id is required", http.StatusBadRequest)
		return
	}

	usage, err := h.Server.GetPowerUsage(serverID, ownerID)
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

package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) GetServerConsole(w http.ResponseWriter, r *http.Request) {
	serverID := strings.TrimSpace(r.URL.Query().Get("server_id"))
	if serverID == "" {
		http.Error(w, "server_id is required", http.StatusBadRequest)
		return
	}

	tail := 200
	if rawTail := strings.TrimSpace(r.URL.Query().Get("tail")); rawTail != "" {
		parsed, err := strconv.Atoi(rawTail)
		if err != nil || parsed <= 0 || parsed > 1000 {
			http.Error(w, "tail must be a number between 1 and 1000", http.StatusBadRequest)
			return
		}
		tail = parsed
	}

	logs, err := h.Server.DockerService.GetConsoleLogs(serverID, tail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(logs))
}

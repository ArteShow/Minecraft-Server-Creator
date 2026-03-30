package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/core"
)

func DeleteBackup(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
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

	if err = core.DeleteBackup(req.ServerID, token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

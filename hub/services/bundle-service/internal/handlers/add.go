package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/bundle-service/internal/repository"
)

func AddBundle(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
    if userID == "" {
        http.Error(w, "userID header missing", http.StatusBadRequest)
        return
    }

	var req AddBundleRequest
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

	if err = repository.AddBundle(userID, req.Bundle); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
} 
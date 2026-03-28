package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/bundle-service/internal/repository"
)

func CreateBundle(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
    if userID == "" {
        http.Error(w, "userID header missing", http.StatusBadRequest)
        return
    }

	var req CreateBundleRequest
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

	key, err := repository.CreateBundleKey(userID, req.Bundle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := CreateBundleResponse{Key: key}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
} 
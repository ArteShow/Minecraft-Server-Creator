package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/host-metadata-service/internal/repository"
)

func GetHostServerMetadata(w http.ResponseWriter, r *http.Request) {
	metadata, err := repository.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := GetHostServerMetadataResponse{Metadata: metadata}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
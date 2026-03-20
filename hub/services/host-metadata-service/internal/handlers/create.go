package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/host-metadata-service/internal/repository"
)

func CreateHostServerMetaDataEntry(w http.ResponseWriter, r *http.Request) {
	id, err := repository.Create([]string{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := CreateHostServerMetadataResponse{HostServerID: id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

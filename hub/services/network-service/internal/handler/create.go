package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/network-service/internal/repository"
)

func CreateHostServer(w http.ResponseWriter, r *http.Request) {
	var req CreateServerRequest
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

	serverID, err := repository.CreateServerFunction(req.IP, req.HostServerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := CreateServerResponse{ServerID: serverID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

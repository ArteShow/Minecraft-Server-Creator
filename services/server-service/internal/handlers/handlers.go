package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service/internal/models"
	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service/internal/server"
)

func CreateServerHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateServerRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := server.CreateServer(req.Version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := models.CreateServerResponse{ServerID: id}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func StartServerHandler(manager server.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.StartServerRequest
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		proc, err := server.StartServer(req.ServerID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		manager.Add(req.ServerID, proc)

		res := models.StartServerResponse{Status: "running"}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
	}
}

func StopServerHandler(manager server.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.StopServerRequest
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		proc, ok := manager.Get(req.ServerID)
		if !ok || proc == nil {
			http.Error(w, "failed to get server process", http.StatusInternalServerError)
			return
		}

		if err := server.StopServer(proc); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		manager.Remove(req.ServerID)

		res := models.StopServerResponse{Status: "stopped"}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
	}
}

func DeleteServerHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteServerRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := server.DeleteServer(req.ServerID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := models.DeleteServerResponse{Status: "deleted"}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

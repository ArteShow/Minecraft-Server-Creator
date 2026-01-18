package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

func (h *Handler) StopServer(w http.ResponseWriter, r *http.Request) {
	var req StopServerRequest
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

	if err = h.Server.StopServer(req.ServerID, req.ContainerID, req.OwnerID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
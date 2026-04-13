package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

type DeletePluginRequest struct {
	ServerID string `json:"server_id"`
	FileName string `json:"file_name"`
}

func (h *Handler) DeletePlugin(w http.ResponseWriter, r *http.Request) {
	var req DeletePluginRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.Server.DeletePlugin(req.ServerID, req.FileName); err != nil {
		http.Error(w, "plugin delete failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

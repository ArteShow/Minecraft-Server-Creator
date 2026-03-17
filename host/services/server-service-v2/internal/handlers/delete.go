package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

func (h *Handler) DeleteServer(w http.ResponseWriter, r *http.Request) {
	ownerID := r.Header.Get("X-Owner-ID")
	if ownerID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req DeleteServerRequest
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

	if err = h.Server.DeleteServer(req.ServerID, ownerID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
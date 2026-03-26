package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/auth-service/internal/repository"
)

func DeleteUsersBundle(w http.ResponseWriter, r *http.Request) {
	var req DeleteUsersBundleRequest
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

	if err = repository.DeleteBundle(req.UserID, req.Bundle); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
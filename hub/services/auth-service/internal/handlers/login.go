package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/auth-service/internal/core"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginUserRequest
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

	token, err := core.LoginUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := LoginUserResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
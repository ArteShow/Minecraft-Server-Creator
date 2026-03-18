package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	client "github.com/ArteShow/Minecraft-Server-Creator/host/services/auth-service/internal/client"
	jwtutil "github.com/ArteShow/Minecraft-Server-Creator/host/services/auth-service/internal/jwt"
	"github.com/ArteShow/Minecraft-Server-Creator/host/services/auth-service/internal/models"
	"github.com/ArteShow/Minecraft-Server-Creator/host/services/auth-service/internal/proto"
)

const JWTTTL = 24 * time.Hour

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req models.RegisterRequest
	err = json.Unmarshal(bytes, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Email == "" || req.Password == "" || req.Username == "" {
		http.Error(w, "invalid username, email or password", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 6 || len(req.Username) < 3 {
		http.Error(w, "username or password are too short", http.StatusBadRequest)
		return
	}

	userClient, err := client.NewUserClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer userClient.Close()

	saveUserRes, err := userClient.SaveUser(&proto.SaveUserRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !saveUserRes.GetSuccess() {
		http.Error(w, "failed to save user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req models.LoginRequest
	err = json.Unmarshal(bytes, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userClient, err := client.NewUserClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer userClient.Close()

	loginUserResponse, err := userClient.LoginUser(&proto.LoginUserRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil || !loginUserResponse.GetOk() || loginUserResponse.GetUserId() == "" {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := jwtutil.GenerateToken(loginUserResponse.GetUserId(), JWTTTL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var res models.LoginResponse
	res.Token = token

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

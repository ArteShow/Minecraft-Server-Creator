package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	client "github.com/ArteShow/Minecraft-Server-Creator/services/auth-service/internal/client"
	jwtutil "github.com/ArteShow/Minecraft-Server-Creator/services/auth-service/internal/jwt"
	"github.com/ArteShow/Minecraft-Server-Creator/services/auth-service/internal/models"
	"github.com/ArteShow/Minecraft-Server-Creator/services/auth-service/internal/proto"
	"github.com/ArteShow/Minecraft-Server-Creator/services/auth-service/pkg/hashing"
)

const JWTTTL = 24*time.Hour

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
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hashedPassword, err := hashing.HashPassword(req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	saveUserRes, err := userClient.SaveUser(proto.SaveUserRequest{
		Username: req.Username,
		Password: string(hashedPassword),
		Email: req.Email,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}else if !saveUserRes.GetSuccess() {
		http.Error(w, "failed to save user", http.StatusInternalServerError)
		return
	}

	var res models.RegisterResponse
	res.ID = saveUserRes.GetId()

	json.NewEncoder(w).Encode(res)
	w.WriteHeader(http.StatusAccepted)
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
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	getUserPasswordRes, err := userClient.GetUserPassword(proto.GetUserPasswordRequest{
		Id: req.ID,
		Username: req.Username,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ok := hashing.ComparePasswords(req.Password, []byte(getUserPasswordRes.GetPassword()))
	if !ok {
		http.Error(w, "no user found with these username and password", http.StatusUnauthorized)
		return
	}

	token, err := jwtutil.GenerateToken(req.ID, JWTTTL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var res models.LoginResponse
	res.Token = token
	json.NewEncoder(w).Encode(res)
	w.WriteHeader(http.StatusAccepted)
}
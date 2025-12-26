package models

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email`
}

type RegisterResponse struct {
	ID string `json:"id"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ID       string `json:"id"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
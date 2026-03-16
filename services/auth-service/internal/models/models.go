package models

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ID       string `json:"id"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

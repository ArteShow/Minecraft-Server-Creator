package handlers

type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Type     string `json:"type"`
	JWT      string `json:"jwt"`
}

type RegisterUserResponse struct {
	UserID string `json:"user_id"`
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}
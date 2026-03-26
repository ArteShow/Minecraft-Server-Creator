package handlers

type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
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

type CreateBundleRequest struct {
	UserID string `json:"user_id"`
	Bundle string `json:"bundle"`
}

type GetUsersBundlesRequest struct {
	UserID string `json:"user_id"`
}

type GetUsersBundlesResponse struct {
	Bundles map[string]int64 `json:"bundles"`
}
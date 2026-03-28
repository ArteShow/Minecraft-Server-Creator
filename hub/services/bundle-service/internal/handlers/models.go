package handlers

type CreateBundleRequest struct {
	UserID string `json:"user_id"`
	Bundle string `json:"bundle"`
}

type CreateBundleResponse struct {
	Key string `json:"key"`
}
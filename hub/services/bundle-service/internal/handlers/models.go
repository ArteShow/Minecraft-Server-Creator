package handlers

type CreateBundleRequest struct {
	Bundle string `json:"bundle"`
}

type CreateBundleResponse struct {
	Key string `json:"key"`
}
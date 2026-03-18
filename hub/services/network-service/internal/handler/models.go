package handler

type CreateServerRequest struct {
	IP string `json:"ip"`
}

type CreateServerResponse struct {
	ServerID string `json:"server_id"`
}

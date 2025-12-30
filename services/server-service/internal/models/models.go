package models

type CreateServerRequest struct {
	Version string `json:"version"`
	UserID  string `json:"id"`
}

type CreateServerResponse struct {
	ServerID string `json:"id"`
}

type StartServerRequest struct {
	ServerID string `json:"server_id"`
	UserID   string `json:"id"`
}

type StartServerResponse struct {
	Status string `json:"status"`
}

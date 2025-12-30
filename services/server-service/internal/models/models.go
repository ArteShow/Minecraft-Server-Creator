package models

type CreateServerRequest struct {
	Version string `json:"version"`
	UserID  string `json:"id"`
}

type CreateServerResponse struct {
	ServerID string `json:"id"`
}

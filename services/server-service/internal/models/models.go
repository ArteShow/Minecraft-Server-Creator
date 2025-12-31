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

type StopServerRequest struct {
	ServerID string `json:"id"`
}

type StopServerResponse struct {
	Status string `json:"status"`
}

type DeleteServerRequest struct {
	ServerID string `json:"id"`
}

type DeleteServerResponse struct {
	Status string `json:"status"`
}

type GetLogRequest struct {
	ServerID string `json:"id"`
}

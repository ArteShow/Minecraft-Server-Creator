package handlers

type CreateServerRequest struct {
	Version string `json:"version"`
}

type CreateServerResponse struct {
	ServerID string `json:"server_id"`
	Port     int    `json:"port"`
}

type StartServerRequest struct {
	ServerID string `json:"server_id"`
}

type StopServerRequest struct {
	ServerID    string `json:"server_id"`
}

type DeleteServerRequest struct {
	ServerID string `json:"server_id"`
}
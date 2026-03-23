package handlers

type CreateServerRequest struct {
	Version string `json:"version"`
}

type CreateServerResponse struct {
	Port     int    `json:"port"`
	ServerID string `json:"server_id"`
}

type StartServerRequest struct {
	ServerID string `json:"server_id"`
}

type StopServerRequest struct {
	ServerID string `json:"server_id"`
}

type DeleteServerRequest struct {
	ServerID string `json:"server_id"`
}

type GetServerStatsRequest struct {
	Key      string `json:"key"`
	ServerID string `json:"server_id"`
}

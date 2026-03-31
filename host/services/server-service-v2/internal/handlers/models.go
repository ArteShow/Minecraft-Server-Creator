package handlers

type CreateServerRequest struct {
	Version string `json:"version"`
	Port    int    `json:"port"`
}

type CreateServerResponse struct {
	ServerID string `json:"server_id"`
	Port     int    `json:"port"`
}

type StartServerRequest struct {
	ServerID string `json:"server_id"`
	RAM      string `json:"RAM"`
	CPUCores int    `json:"cpu_cores"`
}

type StopServerRequest struct {
	ServerID string `json:"server_id"`
}

type DeleteServerRequest struct {
	ServerID string `json:"server_id"`
}

type GetServerStatsRequest struct {
	ServerID string `json:"server_id"`
	Key      string `json:"key"`
}

type GetServerStatsResponse struct {
	Value string `json:"value"`
}

type CreateBackupRequest struct {
	ServerID string `json:"server_id"`
}

type GetBackupRequest struct {
	ServerID string `json:"server_id"`
}

type DeleteBackupRequest struct {
	ServerID string `json:"server_id"`
}

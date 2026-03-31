package handlers

type CreateServerRequest struct {
	Version string `json:"version"`
	Bundle  string `json:"bundle"`
}

type CreateServerResponse struct {
	Port     int    `json:"port"`
	ServerID string `json:"server_id"`
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
	Key      string `json:"key"`
	ServerID string `json:"server_id"`
}

type CreateBackupRequest struct {
	ServerID string `json:"server_id"`
	Bundle   string `json:"bundle"`
}

type GetBackupRequest struct {
	ServerID string `json:"server_id"`
}

type DeleteBackupRequest struct {
	ServerID string `json:"server_id"`
	BackupID string `json:"backup_id"`
}

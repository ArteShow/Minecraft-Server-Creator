package handler

type CreateServerRequest struct {
	IP           string `json:"ip"`
	HostServerID string `json:"host_server_id,omitempty"`
}

type CreateServerResponse struct {
	ServerID string `json:"server_id"`
}

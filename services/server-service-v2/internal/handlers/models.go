package handlers

type CreateServerRequest struct {
	Version string `json:"version"`
	OwnerID string `json:"owner_id"`
}

type CreateServerResponse struct {
	ServerID string `json:"server_id"`
}

type StartServerRequest struct {
	ServerID string `json:"server_id"`
	OwnerID  string `json:"owner_id"`
}

type StartServerResponse struct {
	ContainerID string `json:"container_id"`
}

type StopServerRequest struct {
	ServerID    string `json:"server_id"`
	ContainerID string `json:"container_id"`
	OwnerID     string `json:"owner_id"`
}

type DeleteServerRequest struct {
	ServerID string `json:"server_id"`
	OwnerID  string `json:"owner_id"`
}
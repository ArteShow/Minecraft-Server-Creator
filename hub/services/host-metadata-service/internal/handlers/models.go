package handlers

import "github.com/ArteShow/Minecraft-Server-Creator/hub/services/host-metadata-service/internal/repository"

type CreateHostServerMetadataResponse struct {
	HostServerID string `json:"host_server_id"`
}

type DeleteHostServerMetadataRequest struct {
	HostServerID string `json:"host_server_id"`
}

type GetHostServerMetadataResponse struct {
	Metadata []repository.Host `json:"metadata"`
}

type AddServerToHostServerRequest struct {
	HostServerId string `json:"host_server_id"`
	ServerID     string `json:"server_id"`
}

package handlers

import "github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/internal/core"

type Handler struct {
	Server core.Server
}

type CreateServerRequest struct {
	Version string `json:"version"`
}

type CreateServerResponse struct {
	ServerID string `json:"server_id"`
}

type StartServerRequest struct {
	ServerID string `json:"server_id"`
}

type StartServerResponse struct {
	ContainerID string `json:"container_id"`
}

type StopServerRequest struct {
	ServerID string `json:"server_id"`
	ContainerID string `json:"container_id"`
}

type DeleteServerRequest struct {
	ServerID string `json:"server_id"`
}
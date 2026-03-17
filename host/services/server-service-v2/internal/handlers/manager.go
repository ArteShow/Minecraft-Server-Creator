package handlers

import (
	"sync"

	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/internal/core"
	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/internal/docker"
	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/internal/server"
)

type Handler struct {
	Server core.Server
}

func NewHandler() (*Handler, error) {
	ds, err := docker.NewDockerService()
	if err != nil {
		return nil, err
	}
	return &Handler{
		Server: core.Server{
			DockerService: ds,
			Processes: server.Manager{
				Mu: sync.Mutex{},
				Containers: make(map[string]string),
			},
		},
	}, nil
}
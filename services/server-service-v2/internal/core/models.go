package core

import (
	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/internal/docker"
	"github.com/ArteShow/Minecraft-Server-Creator/services/server-service-v2/internal/server"
)

type Server struct {
	DockerService docker.DockerService
	Processes server.Manager
}
package server

import (
	"context"
	"net/http"
	"time"

	network_pb "github.com/ArteShow/Minecraft-Server-Creator/hub/services/network-service/internal/proto"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/network-service/internal/repository"
)

type Server struct {
	network_pb.UnimplementedNetworkServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetServerStatus(_ context.Context, req *network_pb.GetServerStatusRequest) (*network_pb.GetServerStatusResponse, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get("http://" + req.GetIp())
	if err != nil {
		return &network_pb.GetServerStatusResponse{}, err
	}
	defer resp.Body.Close()

	return &network_pb.GetServerStatusResponse{Status: int64(resp.StatusCode)}, nil
}

func (s *Server) GetServerMetadata(_ context.Context, req *network_pb.GetServerMetadataRequest) (*network_pb.GetServerMetadataResponse, error) {
	metadata, err := repository.GetServerMetadataByID(req.GetServerId())
	if err != nil {
		return &network_pb.GetServerMetadataResponse{}, err
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	start := time.Now()
	resp, err := client.Get("http://" + metadata.IP)
	if err != nil {
		return &network_pb.GetServerMetadataResponse{}, err
	}
	defer resp.Body.Close()

	return &network_pb.GetServerMetadataResponse{Ip: metadata.ID, Ping: int64(time.Since(start))}, nil
}

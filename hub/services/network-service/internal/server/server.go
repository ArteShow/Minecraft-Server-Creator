package server

import (
	"context"

	network_pb "github.com/ArteShow/Minecraft-Server-Creator/hub/services/network-service/internal/proto"
)

type Server struct {
	network_pb.UnimplementedNetworkServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GetServerStatus(_ context.Context, req network_pb.GetServerStatusRequest) (network_pb.GetServerStatusResponse, error) {
	return network_pb.GetServerStatusResponse{}, nil
}

func (s *Server) GetServerMetadata(_ context.Context, req network_pb.GetServerMetadataRequest) (network_pb.GetServerMetadataResponse, error) {
	return network_pb.GetServerMetadataResponse{}, nil
}

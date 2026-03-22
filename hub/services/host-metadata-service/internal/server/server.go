package server

import (
	"context"

	proto_pb "github.com/ArteShow/Minecraft-Server-Creator/hub/services/host-metadata-service/internal/proto"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/host-metadata-service/internal/repository"
)

type Server struct {
	proto_pb.UnimplementedHostMetadataServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) CreateHostServer(_ context.Context, _ *proto_pb.CreateHostServerRequest) (*proto_pb.CreateHostServerResponse, error) {
	id, err := repository.Create([]string{})
	if err != nil {
		return &proto_pb.CreateHostServerResponse{}, err
	}

	return &proto_pb.CreateHostServerResponse{HostServerId: id}, nil
}

func (s *Server) DeleteHostServer(_ context.Context, req *proto_pb.DeleteHostServerRequest) (*proto_pb.DeleteHostServerResponse, error) {
	err := repository.Delete(req.GetHostServerId())
	if err != nil {
		return &proto_pb.DeleteHostServerResponse{}, err
	}

	return &proto_pb.DeleteHostServerResponse{}, nil
}

func (s *Server) GetAllHostServers(_ context.Context, _ *proto_pb.GetAllHostServersRequest) (*proto_pb.GetAllHostServersResponse, error) {
	hosts, err := repository.Get()
	if err != nil {
		return &proto_pb.GetAllHostServersResponse{}, err
	}

	pbHosts := make([]*proto_pb.HostServer, len(hosts))
	for i, host := range hosts {
		pbHosts[i] = &proto_pb.HostServer{
			Id:        host.ID,
			Servers:   host.Servers,
			CreatedAt: host.CreatedAt,
		}
	}

	return &proto_pb.GetAllHostServersResponse{Hosts: pbHosts}, nil
}

func (s *Server) AddServerToHost(_ context.Context, req *proto_pb.AddServerToHostRequest) (*proto_pb.AddServerToHostResponse, error) {
	err := repository.AddServer(req.GetHostServerId(), req.GetServerId())
	if err != nil {
		return &proto_pb.AddServerToHostResponse{}, err
	}

	return &proto_pb.AddServerToHostResponse{}, nil
}

func (s *Server) RemoveServerFromHost(_ context.Context, req *proto_pb.RemoveServerFromHostRequest) (*proto_pb.RemoveServerFromHostResponse, error) {
	err := repository.RemoveServer(req.GetHostServerId(), req.GetServerId())
	if err != nil {
		return &proto_pb.RemoveServerFromHostResponse{}, err
	}

	return &proto_pb.RemoveServerFromHostResponse{}, nil
}

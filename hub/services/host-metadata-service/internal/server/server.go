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

func (s *Server) CreateHostServer(_ context.Context, req *proto_pb.CreateHostServerRequest) (*proto_pb.CreateHostServerResponse, error) {
	id, err := repository.Create(map[string]int{}, req.GetRam(), req.GetCores())
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

	for index, host := range hosts {
		pbServers := make(map[string]int32, len(host.Servers))
		for serverID, port := range host.Servers {
			pbServers[serverID] = int32(port)
		}

		pbHosts[index] = &proto_pb.HostServer{
			Id:        host.ID,
			Servers:   pbServers,
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

func (s *Server) GetRAM(_ context.Context, req *proto_pb.GetRAMRequest) (*proto_pb.GetRAMResponse, error) {
	ram, err := repository.GetRAM(req.GetHostServerId())
	if err != nil {
		return &proto_pb.GetRAMResponse{}, err
	}
	return &proto_pb.GetRAMResponse{Ram: ram}, nil
}

func (s *Server) GetCores(_ context.Context, req *proto_pb.GetCoresRequest) (*proto_pb.GetCoresResponse, error) {
	cores, err := repository.GetCores(req.GetHostServerId())
	if err != nil {
		return &proto_pb.GetCoresResponse{}, err
	}
	return &proto_pb.GetCoresResponse{Cores: cores}, nil
}

func (s *Server) SubtractRAM(_ context.Context, req *proto_pb.SubtractRAMRequest) (*proto_pb.SubtractRAMResponse, error) {
	err := repository.SubtractRAM(req.GetHostServerId(), req.GetRam())
	if err != nil {
		return &proto_pb.SubtractRAMResponse{}, err
	}
	return &proto_pb.SubtractRAMResponse{}, nil
}

func (s *Server) SubtractCores(_ context.Context, req *proto_pb.SubtractCoresRequest) (*proto_pb.SubtractCoresResponse, error) {
	err := repository.SubtractCores(req.GetHostServerId(), req.GetCores())
	if err != nil {
		return &proto_pb.SubtractCoresResponse{}, err
	}
	return &proto_pb.SubtractCoresResponse{}, nil
}

func (s *Server) AddPortToServer(_ context.Context, req *proto_pb.AddPortToServerRequest) (*proto_pb.AddPortToServerResponse, error) {
	err := repository.AddPortToServer(req.GetHostServerId(), req.GetServerId(), int(req.GetPort()))
	if err != nil {
		return &proto_pb.AddPortToServerResponse{}, err
	}
	return &proto_pb.AddPortToServerResponse{}, nil
}

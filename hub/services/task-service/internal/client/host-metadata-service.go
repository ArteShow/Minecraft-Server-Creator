package client

import (
	"context"
	"errors"

	pb "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/host-metadata-service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type HostClient struct {
	Client pb.HostMetadataServiceClient
	Conn   *grpc.ClientConn
}

func NewHostClient() (*HostClient, error) {
	conn, err := grpc.Dial("network-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewHostMetadataServiceClient(conn)
	if client == nil {
		conn.Close()
		return nil, errors.New("failed to create NewHostClient")
	}

	return &HostClient{
		Client: client,
		Conn:   conn,
	}, nil
}

func (u *HostClient) Close() error {
	return u.Conn.Close()
}

func (u *HostClient) CreateHostServer(req *pb.CreateHostServerRequest) (*pb.CreateHostServerResponse, error) {
	res, err := u.Client.CreateHostServer(context.Background(), req)
	if err != nil {
		return &pb.CreateHostServerResponse{}, err
	}

	return res, nil
}

func (u *HostClient) DeleteHostServer(req *pb.DeleteHostServerRequest) (*pb.DeleteHostServerResponse, error) {
	res, err := u.Client.DeleteHostServer(context.Background(), req)
	if err != nil {
		return &pb.DeleteHostServerResponse{}, err
	}

	return res, nil
}

func (u *HostClient) GetAllHostServers(req *pb.GetAllHostServersRequest) (*pb.GetAllHostServersResponse, error) {
	res, err := u.Client.GetAllHostServers(context.Background(), req)
	if err != nil {
		return &pb.GetAllHostServersResponse{}, err
	}

	return res, nil
}

func (u *HostClient) AddServerToHost(req *pb.AddServerToHostRequest) (*pb.AddServerToHostResponse, error) {
	res, err := u.Client.AddServerToHost(context.Background(), req)
	if err != nil {
		return &pb.AddServerToHostResponse{}, err
	}

	return res, nil
}

func (u *HostClient) RemoveServerFromHost(req *pb.RemoveServerFromHostRequest) (*pb.RemoveServerFromHostResponse, error) {
	res, err := u.Client.RemoveServerFromHost(context.Background(), req)
	if err != nil {
		return &pb.RemoveServerFromHostResponse{}, err
	}

	return res, nil
}

func (u *HostClient) GetRAM(req *pb.GetRAMRequest) (*pb.GetRAMResponse, error) {
	res, err := u.Client.GetRAM(context.Background(), req)
	if err != nil {
		return &pb.GetRAMResponse{}, err
	}

	return res, nil
}

func (u *HostClient) GetCores(req *pb.GetCoresRequest) (*pb.GetCoresResponse, error) {
	res, err := u.Client.GetCores(context.Background(), req)
	if err != nil {
		return &pb.GetCoresResponse{}, err
	}

	return res, nil
}

func (u *HostClient) SubtractRAM(req *pb.SubtractRAMRequest) (*pb.SubtractRAMResponse, error) {
	res, err := u.Client.SubtractRAM(context.Background(), req)
	if err != nil {
		return &pb.SubtractRAMResponse{}, err
	}

	return res, nil
}

func (u *HostClient) SubtractCores(req *pb.SubtractCoresRequest) (*pb.SubtractCoresResponse, error) {
	res, err := u.Client.SubtractCores(context.Background(), req)
	if err != nil {
		return &pb.SubtractCoresResponse{}, err
	}

	return res, nil
}

func (u *HostClient) AddPortToServer(req *pb.AddPortToServerRequest) (*pb.AddPortToServerResponse, error) {
	res, err := u.Client.AddPortToServer(context.Background(), req)
	if err != nil {
		return &pb.AddPortToServerResponse{}, err
	}

	return res, nil
}
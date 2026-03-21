package client

import (
	"context"
	"errors"

	pb "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/network-service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NetworkClient struct {
	Client pb.NetworkServiceClient
	Conn   *grpc.ClientConn
}

func NewNetworkClient() (*NetworkClient, error) {
	conn, err := grpc.Dial("network-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewNetworkServiceClient(conn)
	if client == nil {
		conn.Close()
		return nil, errors.New("failed to create NewNetworkClient")
	}

	return &NetworkClient{
		Client: client,
		Conn:   conn,
	}, nil
}

func (u *NetworkClient) Close() error {
	return u.Conn.Close()
}

func (u *NetworkClient) GetServerStatus(req *pb.GetServerStatusRequest) (*pb.GetServerStatusResponse, error) {
	res, err := u.Client.GetServerStatus(context.Background(), req)
	if err != nil {
		return &pb.GetServerStatusResponse{}, err
	}

	return res, nil
}

func (u *NetworkClient) GetServerMetadata(req *pb.GetServerMetadataRequest) (*pb.GetServerMetadataResponse, error) {
	res, err := u.Client.GetServerMetadata(context.Background(), req)
	if err != nil {
		return &pb.GetServerMetadataResponse{}, err
	}

	return res, nil
}

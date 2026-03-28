package client

import (
	"context"
	"errors"

	pb "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/bundle-service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BundleClient struct {
	Client pb.BundleServiceClient
	Conn   *grpc.ClientConn
}

func NewBundleClient() (*BundleClient, error) {
	conn, err := grpc.Dial("bundle-service:50055", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewBundleServiceClient(conn)
	if client == nil {
		conn.Close()
		return nil, errors.New("failed to create NewNetworkClient")
	}

	return &BundleClient{
		Client: client,
		Conn:   conn,
	}, nil
}

func (u *BundleClient) Close() error {
	return u.Conn.Close()
}

func (u *BundleClient) DisableBundleKey(req *pb.DisableBundleKeyRequest) (*pb.DisableBundleKeyResponse, error) {
	res, err := u.Client.DisableBundleKey(context.Background(), req)
	if err != nil {
		return &pb.DisableBundleKeyResponse{}, err
	}

	return res, nil
}

func (u *BundleClient) AddBundle(req *pb.AddBundleRequest) (*pb.AddBundleResponse, error) {
	res, err := u.Client.AddBundle(context.Background(), req)
	if err != nil {
		return &pb.AddBundleResponse{}, err
	}

	return res, nil
}

func (u *BundleClient) RemoveBundle(req *pb.RemoveBundleRequest) (*pb.RemoveBundleResponse, error) {
	res, err := u.Client.RemoveBundle(context.Background(), req)
	if err != nil {
		return &pb.RemoveBundleResponse{}, err
	}

	return res, nil
}
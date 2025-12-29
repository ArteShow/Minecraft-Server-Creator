package client

import (
	"context"
	"errors"

	pb "github.com/ArteShow/Minecraft-Server-Creator/services/auth-service/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	Client pb.UserServiceClient
}

func NewUserClient() (*UserClient, error) {
	conn, err := grpc.NewClient("user-service:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			return
		}
	}()

	client := pb.NewUserServiceClient(conn)
	if client == nil {
		return nil, errors.New("failed to create NewUserClient")
	}

	return &UserClient{Client: client}, nil
}

func (u *UserClient) SaveUser(req *pb.SaveUserRequest) (*pb.SaveUserResponse, error) {
	res, err := u.Client.SaveUser(context.Background(), req)
	if err != nil {
		return &pb.SaveUserResponse{}, err
	}

	return res, nil
}

func (u *UserClient) GetUserPassword(req *pb.GetUserPasswordRequest) (*pb.GetUserPasswordResponse, error) {
	res, err := u.Client.GetUserPassword(context.Background(), req)
	if err != nil {
		return &pb.GetUserPasswordResponse{}, err
	}

	return res, nil
}

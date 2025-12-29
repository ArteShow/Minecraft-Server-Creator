package server

import (
	"context"

	user_pb "github.com/ArteShow/Minecraft-Server-Creator/user-service/internal/proto"
	"github.com/ArteShow/Minecraft-Server-Creator/user-service/internal/repository"
)

type Server struct {
	user_pb.UnimplementedUserServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) SaveUser(_ context.Context, req *user_pb.SaveUserRequest) (*user_pb.SaveUserResponse, error) {
	id, err := repository.CreateUser(req.GetUsername(), req.GetPassword(), req.GetEmail())
	if err != nil {
		return &user_pb.SaveUserResponse{}, err
	}

	return &user_pb.SaveUserResponse{
		Id:      id,
		Success: true,
	}, nil
}

func (s *Server) GetUserPassword(_ context.Context, req *user_pb.GetUserPasswordRequest) (*user_pb.GetUserPasswordResponse, error) {
	password, err := repository.GetPassword(req.GetId())
	if err != nil {
		return &user_pb.GetUserPasswordResponse{}, err
	}

	return &user_pb.GetUserPasswordResponse{
		Password: password,
	}, nil
}

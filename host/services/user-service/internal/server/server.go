package server

import (
	"context"

	user_pb "github.com/ArteShow/Minecraft-Server-Creator/host/services/user-service/internal/proto"
	"github.com/ArteShow/Minecraft-Server-Creator/host/services/user-service/internal/repository"
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

func (s *Server) LoginUser(_ context.Context, req *user_pb.LoginUserRequest) (*user_pb.LoginUserResponse, error) {
	id, err := repository.LoginUser(req.GetUsername(), req.GetPassword())
	if err != nil {
		return &user_pb.LoginUserResponse{Ok: false, UserId: ""}, err
	}

	return &user_pb.LoginUserResponse{
		UserId: id,
		Ok:     true,
	}, nil
}

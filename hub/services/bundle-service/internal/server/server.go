package server

import (
	"context"

	proto_pb "github.com/ArteShow/Minecraft-Server-Creator/hub/services/bundle-service/internal/proto"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/bundle-service/internal/repository"
)

type Server struct {
	proto_pb.UnimplementedBundleServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) DisableBundleKey(_ context.Context, req *proto_pb.DisableBundleKeyRequest) (*proto_pb.DisableBundleKeyResponse, error) {
	ok, _, _, err := repository.UseBundleKey(req.GetKey())
	if err != nil || !ok{
		return &proto_pb.DisableBundleKeyResponse{}, err
	}

	return &proto_pb.DisableBundleKeyResponse{}, nil
}

func (s *Server) AddBundle(_ context.Context, req *proto_pb.AddBundleRequest) (*proto_pb.AddBundleResponse, error) {
	if err := repository.AddBundle(req.GetUserID(), req.GetBundle()); err != nil{
		return &proto_pb.AddBundleResponse{}, err
	}

	return &proto_pb.AddBundleResponse{}, nil
}

func (s *Server) RemoveBundle(_ context.Context, req *proto_pb.RemoveBundleRequest) (*proto_pb.RemoveBundleResponse, error) {
	if err := repository.DeleteBundle(req.GetUserID(), req.GetBundle()); err != nil{
		return &proto_pb.RemoveBundleResponse{}, err
	}

	return &proto_pb.RemoveBundleResponse{}, nil
}
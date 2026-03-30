package server

import (
	"context"

	proto_pb "github.com/ArteShow/Minecraft-Server-Creator/hub/services/backup-service/internal/proto"
	"github.com/ArteShow/Minecraft-Server-Creator/hub/services/backup-service/internal/repository"
)

type Server struct {
	proto_pb.UnimplementedBundleServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) CreateBackup(_ context.Context, req *proto_pb.CreateBackupRequest) (*proto_pb.CreateBackupResponse, error) {
	if err := repository.CreateBackup(context.Background(), &repository.Backup{
		ServerID: req.GetServerID(),
		UserID: req.GetUserID(),
	}); err != nil {
		return &proto_pb.CreateBackupResponse{}, err
	}

	return &proto_pb.CreateBackupResponse{}, nil
}

func (s *Server) GetBackup(_ context.Context, req *proto_pb.GetBackupRequest) (*proto_pb.GetBackupResponse, error) {
	backups, err := repository.GetBackups(context.Background(), req.GetServerID())
	if err != nil {
		return &proto_pb.GetBackupResponse{}, err
	}

	var resp []*proto_pb.BackupList
	for _, backup := range backups {
		resp = append(resp, &proto_pb.BackupList{ 
			ServerID: backup.ServerID,
			Backup:   backup.BackupID,
			UserID:   backup.UserID,
		})
	}

	return &proto_pb.GetBackupResponse{Backups: resp}, nil
}

func (s *Server) DeleteBackup(_ context.Context, req *proto_pb.DeleteBackupRequest) (*proto_pb.DeleteBackupResponse, error) {
	if err := repository.DeleteBackup(context.Background(), req.GetBackupID()); err != nil{
		return &proto_pb.DeleteBackupResponse{}, err
	}

	return &proto_pb.DeleteBackupResponse{}, nil
}
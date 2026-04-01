package client

import (
	"context"
	"errors"

	pb "github.com/ArteShow/Minecraft-Server-Creator/hub/services/task-service/internal/proto/backup-service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BackupClient struct {
	Client pb.BackupServiceClient
	Conn   *grpc.ClientConn
}

func NewBackupClient() (*BackupClient, error) {
	conn, err := grpc.Dial("backup-service:50056", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewBackupServiceClient(conn)
	if client == nil {
		conn.Close()
		return nil, errors.New("failed to create NewBackupClient")
	}

	return &BackupClient{
		Client: client,
		Conn:   conn,
	}, nil
}

func (u *BackupClient) Close() error {
	return u.Conn.Close()
}

func (u *BackupClient) CreateBackup(req *pb.CreateBackupRequest) (*pb.CreateBackupResponse, error) {
	res, err := u.Client.CreateBackup(context.Background(), req)
	if err != nil {
		return &pb.CreateBackupResponse{}, err
	}

	return res, nil
}

func (u *BackupClient) DeleteBackup(req *pb.DeleteBackupRequest) (*pb.DeleteBackupResponse, error) {
	res, err := u.Client.DeleteBackup(context.Background(), req)
	if err != nil {
		return &pb.DeleteBackupResponse{}, err
	}

	return res, nil
}

func (u *BackupClient) GetBackup(req *pb.GetBackupRequest) (*pb.GetBackupResponse, error) {
	res, err := u.Client.GetBackup(context.Background(), req)
	if err != nil {
		return &pb.GetBackupResponse{}, err
	}

	return res, nil
}